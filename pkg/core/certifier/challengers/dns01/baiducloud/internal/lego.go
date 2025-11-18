package internal

import (
	"errors"
	"fmt"
	"time"

	bce "github.com/baidubce/bce-sdk-go/bce"
	bcedns "github.com/baidubce/bce-sdk-go/services/dns"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/platform/config/env"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/samber/lo"
)

const (
	envNamespace = "BAIDUCLOUD_"

	EnvAccessKeyID     = envNamespace + "ACCESS_KEY_ID"
	EnvSecretAccessKey = envNamespace + "SECRET_ACCESS_KEY"

	EnvTTL                = envNamespace + "TTL"
	EnvPropagationTimeout = envNamespace + "PROPAGATION_TIMEOUT"
	EnvPollingInterval    = envNamespace + "POLLING_INTERVAL"
	EnvHTTPTimeout        = envNamespace + "HTTP_TIMEOUT"
)

var _ challenge.ProviderTimeout = (*DNSProvider)(nil)

type Config struct {
	AccessKeyID     string
	SecretAccessKey string

	PropagationTimeout time.Duration
	PollingInterval    time.Duration
	TTL                int
	HTTPTimeout        time.Duration
}

type DNSProvider struct {
	client *bcedns.Client
	config *Config
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
	values, err := env.Get(EnvAccessKeyID, EnvSecretAccessKey)
	if err != nil {
		return nil, fmt.Errorf("baiducloud: %w", err)
	}

	config := NewDefaultConfig()
	config.AccessKeyID = values[EnvAccessKeyID]
	config.SecretAccessKey = values[EnvSecretAccessKey]

	return NewDNSProviderConfig(config)
}

func NewDNSProviderConfig(config *Config) (*DNSProvider, error) {
	if config == nil {
		return nil, errors.New("baiducloud: the configuration of the DNS provider is nil")
	}

	client, err := bcedns.NewClient(config.AccessKeyID, config.SecretAccessKey, "")
	if err != nil {
		return nil, err
	} else {
		if client.Config == nil {
			client.Config = &bce.BceClientConfiguration{}
		}
		client.Config.HTTPClientTimeout = &config.HTTPTimeout
	}

	return &DNSProvider{
		client: client,
		config: config,
	}, nil
}

func (d *DNSProvider) Present(domain, token, keyAuth string) error {
	info := dns01.GetChallengeInfo(domain, keyAuth)

	authZone, err := dns01.FindZoneByFqdn(info.EffectiveFQDN)
	if err != nil {
		return fmt.Errorf("baiducloud: could not find zone for domain %q: %w", domain, err)
	}

	subDomain, err := dns01.ExtractSubDomain(info.EffectiveFQDN, authZone)
	if err != nil {
		return fmt.Errorf("baiducloud: %w", err)
	}

	// REF: https://cloud.baidu.com/doc/DNS/s/El4s7lssr#%E6%B7%BB%E5%8A%A0%E8%A7%A3%E6%9E%90%E8%AE%B0%E5%BD%95
	bceCreateRecordReq := &bcedns.CreateRecordRequest{
		Type:        "TXT",
		Rr:          subDomain,
		Value:       info.Value,
		Description: lo.ToPtr("certimate acme"),
		Ttl:         lo.ToPtr(int32(d.config.TTL)),
	}
	if err := d.client.CreateRecord(dns01.UnFqdn(authZone), bceCreateRecordReq, security.RandomString(32)); err != nil {
		return fmt.Errorf("baiducloud: error when create record: %w", err)
	}

	return nil
}

func (d *DNSProvider) CleanUp(domain, token, keyAuth string) error {
	info := dns01.GetChallengeInfo(domain, keyAuth)

	authZone, err := dns01.FindZoneByFqdn(info.EffectiveFQDN)
	if err != nil {
		return fmt.Errorf("baiducloud: could not find zone for domain %q: %w", domain, err)
	}

	subDomain, err := dns01.ExtractSubDomain(info.EffectiveFQDN, authZone)
	if err != nil {
		return fmt.Errorf("baiducloud: %w", err)
	}

	record, err := d.findDNSRecord(dns01.UnFqdn(authZone), subDomain, info.Value)
	if err != nil {
		return fmt.Errorf("baiducloud: error when find record: %q: %w", domain, err)
	}

	// REF: https://cloud.baidu.com/doc/DNS/s/El4s7lssr#%E5%88%A0%E9%99%A4%E8%A7%A3%E6%9E%90%E8%AE%B0%E5%BD%95
	if err := d.client.DeleteRecord(dns01.UnFqdn(authZone), record.Id, security.RandomString(32)); err != nil {
		return fmt.Errorf("baiducloud: %w", err)
	}

	return nil
}

func (d *DNSProvider) Timeout() (timeout, interval time.Duration) {
	return d.config.PropagationTimeout, d.config.PollingInterval
}

func (d *DNSProvider) findDNSRecord(zoneName, subDomain, tokenValue string) (*bcedns.Record, error) {
	bceListRecordPageMarker := ""
	for {
		// REF: https://cloud.baidu.com/doc/DNS/s/El4s7lssr#%E6%9F%A5%E8%AF%A2%E8%A7%A3%E6%9E%90%E8%AE%B0%E5%BD%95%E5%88%97%E8%A1%A8
		bceListRecordReq := &bcedns.ListRecordRequest{}
		bceListRecordReq.Rr = subDomain
		bceListRecordReq.Marker = bceListRecordPageMarker
		bceListRecordReq.MaxKeys = 1000

		bceListRecordResp, err := d.client.ListRecord(zoneName, bceListRecordReq)
		if err != nil {
			return nil, err
		}

		for _, record := range bceListRecordResp.Records {
			if record.Type == "TXT" && record.Rr == subDomain && record.Value == tokenValue {
				return &record, nil
			}
		}

		if bceListRecordResp.NextMarker == "" {
			break
		}

		bceListRecordPageMarker = bceListRecordResp.NextMarker
	}

	return nil, errors.New("could not find record")
}
