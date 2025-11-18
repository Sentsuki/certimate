package internal

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	qs "github.com/google/go-querystring/query"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// This is a modified fork of https://github.com/go-acme/lego/blob/master/providers/dns/westcn/internal/client.go
// It only changes the base URL.
type DnsClient struct {
	username string
	password string

	encoder *encoding.Encoder

	baseURL    *url.URL
	HTTPClient *http.Client
}

func NewDnsClient(username, password string) (*DnsClient, error) {
	if username == "" || password == "" {
		return nil, errors.New("35cn: credentials missing")
	}

	baseURL, _ := url.Parse("https://api.35.cn/api/v2")

	return &DnsClient{
		username:   username,
		password:   password,
		encoder:    simplifiedchinese.GBK.NewEncoder(),
		baseURL:    baseURL,
		HTTPClient: &http.Client{},
	}, nil
}

type DnsAPIResponse[T any] struct {
	Result    int    `json:"result,omitempty"`
	ClientID  string `json:"clientid,omitempty"`
	Message   string `json:"msg,omitempty"`
	ErrorCode int    `json:"errcode,omitempty"`
	Data      T      `json:"data,omitempty"`
}

func (a DnsAPIResponse[T]) Error() string {
	return fmt.Sprintf("%d: %s (%d)", a.ErrorCode, a.Message, a.Result)
}

type DnsRecord struct {
	Domain   string `url:"domain,omitempty"`
	Host     string `url:"host,omitempty"`
	Type     string `url:"type,omitempty"`
	Value    string `url:"value,omitempty"`
	TTL      int    `url:"ttl,omitempty"` // 60~86400 seconds
	Priority int    `url:"level,omitempty"`
}

type DnsRecordID struct {
	ID int `json:"id,omitempty"`
}

func (c *DnsClient) AddRecord(ctx context.Context, record DnsRecord) (int, error) {
	values, err := qs.Values(record)
	if err != nil {
		return 0, err
	}

	req, err := c.newRequest(ctx, "domain", "adddnsrecord", values)
	if err != nil {
		return 0, err
	}

	results := &DnsAPIResponse[DnsRecordID]{}

	err = c.doRequest(req, results)
	if err != nil {
		return 0, err
	}

	if results.Result != http.StatusOK {
		return 0, results
	}

	return results.Data.ID, nil
}

func (c *DnsClient) DeleteRecord(ctx context.Context, domain string, recordID int) error {
	values := url.Values{}
	values.Set("domain", domain)
	values.Set("id", strconv.Itoa(recordID))

	req, err := c.newRequest(ctx, "domain", "deldnsrecord", values)
	if err != nil {
		return err
	}

	results := &DnsAPIResponse[any]{}

	err = c.doRequest(req, results)
	if err != nil {
		return err
	}

	if results.Result != http.StatusOK {
		return results
	}

	return nil
}

func (c *DnsClient) newRequest(ctx context.Context, p, act string, form url.Values) (*http.Request, error) {
	if form == nil {
		form = url.Values{}
	}

	c.sign(form, time.Now())

	values, err := c.convertURLValues(form)
	if err != nil {
		return nil, err
	}

	endpoint := c.baseURL.JoinPath(p, "/")

	query := endpoint.Query()
	query.Set("act", act)
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

func (c *DnsClient) doRequest(req *http.Request, result any) error {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return c.parseHttpError(req, resp)
	}

	if result == nil {
		return nil
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("35cn: %w", err)
	}

	err = c.newGBKDecoder(raw).Decode(result)
	if err != nil {
		return fmt.Errorf("35cn: %w", err)
	}

	return nil
}

func (c *DnsClient) sign(form url.Values, now time.Time) {
	timestamp := strconv.FormatInt(now.UnixMilli(), 10)

	sum := md5.Sum([]byte(c.username + c.password + timestamp))

	form.Set("token", hex.EncodeToString(sum[:]))
	form.Set("username", c.username)
	form.Set("time", timestamp)
}

func (c *DnsClient) convertURLValues(values url.Values) (url.Values, error) {
	results := make(url.Values)

	for key, vs := range values {
		encKey, err := c.encoder.String(key)
		if err != nil {
			return nil, err
		}

		for _, value := range vs {
			encValue, err := c.encoder.String(value)
			if err != nil {
				return nil, err
			}

			results.Add(encKey, encValue)
		}
	}

	return results, nil
}

func (c *DnsClient) parseHttpError(req *http.Request, resp *http.Response) error {
	raw, _ := io.ReadAll(resp.Body)

	result := &DnsAPIResponse[any]{}

	err := c.newGBKDecoder(raw).Decode(result)
	if err != nil {
		return fmt.Errorf("35cn: %w", err)
	}

	return result
}

func (c *DnsClient) newGBKDecoder(raw []byte) *json.Decoder {
	return json.NewDecoder(transform.NewReader(bytes.NewBuffer(raw), simplifiedchinese.GBK.NewDecoder()))
}
