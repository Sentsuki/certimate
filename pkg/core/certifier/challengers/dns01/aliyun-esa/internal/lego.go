package internal

import (
	"errors"
	"fmt"
	"sync"
	"time"

	aliopen "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	aliesa "github.com/alibabacloud-go/esa-20240910/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/platform/config/env"
)

const (
	envNamespace = "ALICLOUDESA_"

	EnvAccessKey = envNamespace + "ACCESS_KEY"
	EnvSecretKey = envNamespace + "SECRET_KEY"
	EnvRegionID  = envNamespace + "REGION_ID"

	EnvTTL                = envNamespace + "TTL"
	EnvPropagationTimeout = envNamespace + "PROPAGATION_TIMEOUT"
	EnvPollingInterval    = envNamespace + "POLLING_INTERVAL"
	EnvHTTPTimeout        = envNamespace + "HTTP_TIMEOUT"
)

var _ challenge.ProviderTimeout = (*DNSProvider)(nil)

type Config struct {
	SecretID  string
	SecretKey string
	RegionID  string

	PropagationTimeout time.Duration
	PollingInterval    time.Duration
	TTL                int
	HTTPTimeout        time.Duration
}

type DNSProvider struct {
	client *EsaClient
	config *Config

	recordIDs   map[string]int64
	recordIDsMu sync.Mutex
}

func NewDefaultConfig() *Config {
	return &Config{
		TTL:                env.GetOrDefaultInt(EnvTTL, 300),
		PropagationTimeout: env.GetOrDefaultSecond(EnvPropagationTimeout, 5*time.Minute),
		PollingInterval:    env.GetOrDefaultSecond(EnvPollingInterval, dns01.DefaultPollingInterval),
		HTTPTimeout:        env.GetOrDefaultSecond(EnvHTTPTimeout, 30*time.Second),
	}
}

func NewDNSProvider() (*DNSProvider, error) {
	values, err := env.Get(EnvAccessKey, EnvSecretKey, EnvRegionID)
	if err != nil {
		return nil, fmt.Errorf("alicloud-esa: %w", err)
	}

	config := NewDefaultConfig()
	config.SecretID = values[EnvAccessKey]
	config.SecretKey = values[EnvSecretKey]
	config.RegionID = values[EnvRegionID]

	return NewDNSProviderConfig(config)
}

func NewDNSProviderConfig(config *Config) (*DNSProvider, error) {
	if config == nil {
		return nil, errors.New("alicloud-esa: the configuration of the DNS provider is nil")
	}

	if config.RegionID == "" {
		config.RegionID = "cn-hangzhou"
	}

	client, err := NewEsaClient(&aliopen.Config{
		AccessKeyId:     tea.String(config.SecretID),
		AccessKeySecret: tea.String(config.SecretKey),
		Endpoint:        tea.String(fmt.Sprintf("esa.%s.aliyuncs.com", config.RegionID)),
	})
	if err != nil {
		return nil, fmt.Errorf("alicloud-esa: %w", err)
	}

	return &DNSProvider{
		client:      client,
		config:      config,
		recordIDs:   make(map[string]int64),
		recordIDsMu: sync.Mutex{},
	}, nil
}

func (d *DNSProvider) Present(domain, token, keyAuth string) error {
	info := dns01.GetChallengeInfo(domain, keyAuth)

	authZone, err := dns01.FindZoneByFqdn(info.EffectiveFQDN)
	if err != nil {
		return fmt.Errorf("alicloud-esa: could not find zone for domain %q: %w", domain, err)
	}

	siteName := dns01.UnFqdn(authZone)
	siteID, err := d.findSiteIdByName(siteName)
	if err != nil {
		return fmt.Errorf("alicloud-esa: could not find site for zone %q: %w", siteName, err)
	}

	// REF: https://www.alibabacloud.com/help/en/edge-security-acceleration/esa/api-esa-2024-09-10-createrecord
	aliCreateRecordReq := &aliesa.CreateRecordRequest{
		SiteId:     tea.Int64(siteID),
		Type:       tea.String("TXT"),
		RecordName: tea.String(dns01.UnFqdn(info.EffectiveFQDN)),
		Data: &aliesa.CreateRecordRequestData{
			Value: tea.String(info.Value),
		},
		Ttl: tea.Int32(int32(d.config.TTL)),
	}
	aliCreateRecordResp, err := d.client.CreateRecord(aliCreateRecordReq)
	if err != nil {
		return fmt.Errorf("alicloud-esa: error when create record: %w", err)
	}

	d.recordIDsMu.Lock()
	d.recordIDs[token] = *aliCreateRecordResp.Body.GetRecordId()
	d.recordIDsMu.Unlock()

	return nil
}

func (d *DNSProvider) CleanUp(domain, token, keyAuth string) error {
	info := dns01.GetChallengeInfo(domain, keyAuth)

	d.recordIDsMu.Lock()
	recordID, ok := d.recordIDs[token]
	d.recordIDsMu.Unlock()
	if !ok {
		return fmt.Errorf("alicloud-esa: unknown record ID for '%s'", info.EffectiveFQDN)
	}

	// REF: https://www.alibabacloud.com/help/en/edge-security-acceleration/esa/api-esa-2024-09-10-deleterecord
	aliDeleteRecordReq := &aliesa.DeleteRecordRequest{
		RecordId: &recordID,
	}
	if _, err := d.client.DeleteRecord(aliDeleteRecordReq); err != nil {
		return fmt.Errorf("alicloud-esa: error when delete record %w", err)
	}

	return nil
}

func (d *DNSProvider) Timeout() (timeout, interval time.Duration) {
	return d.config.PropagationTimeout, d.config.PollingInterval
}

func (d *DNSProvider) findSiteIdByName(siteName string) (int64, error) {
	aliListSitesPageNumber := 1
	aliListSitesPageSize := 500
	for {
		// REF: https://www.alibabacloud.com/help/en/edge-security-acceleration/esa/api-esa-2024-09-10-listsites
		aliListSitesReq := &aliesa.ListSitesRequest{
			SiteName:       tea.String(siteName),
			SiteSearchType: tea.String("exact"),
			AccessType:     tea.String("NS"),
			PageNumber:     tea.Int32(int32(aliListSitesPageNumber)),
			PageSize:       tea.Int32(int32(aliListSitesPageSize)),
		}
		aliListSitesResp, err := d.client.ListSites(aliListSitesReq)
		if err != nil {
			return 0, err
		}

		if aliListSitesResp.Body == nil {
			break
		}

		for _, siteItem := range aliListSitesResp.Body.Sites {
			if *siteItem.GetSiteName() == siteName {
				return *siteItem.GetSiteId(), nil
			}
		}

		if len(aliListSitesResp.Body.Sites) < aliListSitesPageSize {
			break
		}

		aliListSitesPageNumber++
	}

	return 0, fmt.Errorf("could not find site '%s'", siteName)
}
