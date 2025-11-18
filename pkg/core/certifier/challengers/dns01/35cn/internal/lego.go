package internal

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/platform/config/env"
)

const (
	envNamespace = "35CN_"

	EnvUsername = envNamespace + "USERNAME"
	EnvPassword = envNamespace + "PASSWORD"

	EnvTTL                = envNamespace + "TTL"
	EnvPropagationTimeout = envNamespace + "PROPAGATION_TIMEOUT"
	EnvPollingInterval    = envNamespace + "POLLING_INTERVAL"
	EnvHTTPTimeout        = envNamespace + "HTTP_TIMEOUT"
)

var _ challenge.ProviderTimeout = (*DNSProvider)(nil)

type Config struct {
	Username string
	Password string

	PropagationTimeout time.Duration
	PollingInterval    time.Duration
	TTL                int
	HTTPClient         *http.Client
}

func NewDefaultConfig() *Config {
	return &Config{
		TTL:                env.GetOrDefaultInt(EnvTTL, 60),
		PropagationTimeout: env.GetOrDefaultSecond(EnvPropagationTimeout, dns01.DefaultPropagationTimeout),
		PollingInterval:    env.GetOrDefaultSecond(EnvPollingInterval, dns01.DefaultPollingInterval),
		HTTPClient: &http.Client{
			Timeout: env.GetOrDefaultSecond(EnvHTTPTimeout, 30*time.Second),
		},
	}
}

type DNSProvider struct {
	config *Config
	client *DnsClient

	recordIDs   map[string]int
	recordIDsMu sync.Mutex
}

func NewDNSProvider() (*DNSProvider, error) {
	values, err := env.Get(EnvUsername, EnvPassword)
	if err != nil {
		return nil, fmt.Errorf("35cn: %w", err)
	}

	config := NewDefaultConfig()
	config.Username = values[EnvUsername]
	config.Password = values[EnvPassword]

	return NewDNSProviderConfig(config)
}

func NewDNSProviderConfig(config *Config) (*DNSProvider, error) {
	if config == nil {
		return nil, errors.New("35cn: the configuration of the DNS provider is nil")
	}

	client, err := NewDnsClient(config.Username, config.Password)
	if err != nil {
		return nil, fmt.Errorf("35cn: %w", err)
	}

	if config.HTTPClient != nil {
		client.HTTPClient = config.HTTPClient
	}

	return &DNSProvider{
		config:    config,
		client:    client,
		recordIDs: make(map[string]int),
	}, nil
}

func (d *DNSProvider) Present(domain, token, keyAuth string) error {
	info := dns01.GetChallengeInfo(domain, keyAuth)

	authZone, err := dns01.FindZoneByFqdn(info.EffectiveFQDN)
	if err != nil {
		return fmt.Errorf("35cn: could not find zone for domain %q: %w", domain, err)
	}

	subDomain, err := dns01.ExtractSubDomain(info.EffectiveFQDN, authZone)
	if err != nil {
		return fmt.Errorf("35cn: %w", err)
	}

	record := DnsRecord{
		Domain: dns01.UnFqdn(authZone),
		Host:   subDomain,
		Type:   "TXT",
		Value:  info.Value,
		TTL:    d.config.TTL,
	}

	recordID, err := d.client.AddRecord(context.Background(), record)
	if err != nil {
		return fmt.Errorf("35cn: error when create record: %w", err)
	}

	d.recordIDsMu.Lock()
	d.recordIDs[token] = recordID
	d.recordIDsMu.Unlock()

	return nil
}

func (d *DNSProvider) CleanUp(domain, token, keyAuth string) error {
	info := dns01.GetChallengeInfo(domain, keyAuth)

	authZone, err := dns01.FindZoneByFqdn(info.EffectiveFQDN)
	if err != nil {
		return fmt.Errorf("35cn: could not find zone for domain %q: %w", domain, err)
	}

	d.recordIDsMu.Lock()
	recordID, ok := d.recordIDs[token]
	d.recordIDsMu.Unlock()

	if !ok {
		return fmt.Errorf("35cn: unknown record ID for '%s' '%s'", info.EffectiveFQDN, token)
	}

	err = d.client.DeleteRecord(context.Background(), dns01.UnFqdn(authZone), recordID)
	if err != nil {
		return fmt.Errorf("35cn: error when delete record: %w", err)
	}

	d.recordIDsMu.Lock()
	delete(d.recordIDs, token)
	d.recordIDsMu.Unlock()

	return nil
}

func (d *DNSProvider) Timeout() (timeout, interval time.Duration) {
	return d.config.PropagationTimeout, d.config.PollingInterval
}
