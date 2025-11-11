package volcenginecdn

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	vecdn "github.com/volcengine/volcengine-go-sdk/service/cdn"
	ve "github.com/volcengine/volcengine-go-sdk/volcengine"
	vesession "github.com/volcengine/volcengine-go-sdk/volcengine/session"

	"github.com/certimate-go/certimate/pkg/core"
	sslmgrsp "github.com/certimate-go/certimate/pkg/core/ssl-manager/providers/volcengine-cdn"
	xcerthostname "github.com/certimate-go/certimate/pkg/utils/cert/hostname"
)

type SSLDeployerProviderConfig struct {
	// 火山引擎 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 火山引擎 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 域名匹配模式。
	// 零值时默认值 [MATCH_PATTERN_EXACT]。
	MatchPattern string `json:"matchPattern,omitempty"`
	// 加速域名（支持泛域名）。
	Domain string `json:"domain"`
}

type SSLDeployerProvider struct {
	config     *SSLDeployerProviderConfig
	logger     *slog.Logger
	sdkClient  vecdn.CDNAPI
	sslManager core.SSLManager
}

var _ core.SSLDeployer = (*SSLDeployerProvider)(nil)

func NewSSLDeployerProvider(config *SSLDeployerProviderConfig) (*SSLDeployerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	client, err := createSDKClient(config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("could not create sdk client: %w", err)
	}

	sslmgr, err := sslmgrsp.NewSSLManagerProvider(&sslmgrsp.SSLManagerProviderConfig{
		AccessKeyId:     config.AccessKeyId,
		AccessKeySecret: config.AccessKeySecret,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create ssl manager: %w", err)
	}

	return &SSLDeployerProvider{
		config:     config,
		logger:     slog.Default(),
		sdkClient:  client,
		sslManager: sslmgr,
	}, nil
}

func (d *SSLDeployerProvider) SetLogger(logger *slog.Logger) {
	if logger == nil {
		d.logger = slog.New(slog.DiscardHandler)
	} else {
		d.logger = logger
	}

	d.sslManager.SetLogger(logger)
}

func (d *SSLDeployerProvider) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*core.SSLDeployResult, error) {
	// 上传证书
	upres, err := d.sslManager.Upload(ctx, certPEM, privkeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to upload certificate file: %w", err)
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	// 获取待部署的 CDN 实例
	domains := make([]string, 0)
	switch d.config.MatchPattern {
	case "", MATCH_PATTERN_EXACT:
		{
			if d.config.Domain == "" {
				return nil, errors.New("config `domain` is required")
			}

			domains = append(domains, d.config.Domain)
		}

	case MATCH_PATTERN_WILDCARD:
		{
			if d.config.Domain == "" {
				return nil, errors.New("config `domain` is required")
			}

			if strings.HasPrefix(d.config.Domain, "*.") {
				temp, err := d.getMatchedDomainsByWildcard(ctx, d.config.Domain)
				if err != nil {
					return nil, err
				}

				domains = temp
			} else {
				domains = append(domains, d.config.Domain)
			}
		}

	case MATCH_PATTERN_CERTSAN:
		{
			temp, err := d.getMatchedDomainsByCertId(ctx, upres.CertId)
			if err != nil {
				return nil, err
			}

			domains = temp
		}

	default:
		return nil, fmt.Errorf("unsupported match pattern: '%s'", d.config.MatchPattern)
	}

	// 遍历绑定证书
	if len(domains) == 0 {
		d.logger.Info("no cdn domains to deploy")
	} else {
		d.logger.Info("found cdn domains to deploy", slog.Any("domains", domains))
		var errs []error

		for _, domain := range domains {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				if err := d.bindCert(ctx, domain, upres.CertId); err != nil {
					errs = append(errs, err)
				}
			}
		}

		if len(errs) > 0 {
			return nil, errors.Join(errs...)
		}
	}

	return &core.SSLDeployResult{}, nil
}

func (d *SSLDeployerProvider) getMatchedDomainsByWildcard(ctx context.Context, wildcardDomain string) ([]string, error) {
	domains := make([]string, 0)

	// 遍历获取加速域名列表，获取匹配的域名
	// REF: https://www.volcengine.com/docs/6454/75269
	listCdnDomainsPageNum := int64(1)
	listCdnDomainsPageSize := int64(100)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		listCdnDomainsReq := &vecdn.ListCdnDomainsInput{
			Domain:   ve.String(strings.TrimPrefix(wildcardDomain, "*.")),
			Status:   ve.String("online"),
			PageNum:  ve.Int64(listCdnDomainsPageNum),
			PageSize: ve.Int64(listCdnDomainsPageSize),
		}
		listCdnDomainsResp, err := d.sdkClient.ListCdnDomains(listCdnDomainsReq)
		d.logger.Debug("sdk request 'cdn.ListCdnDomains'", slog.Any("request", listCdnDomainsReq), slog.Any("response", listCdnDomainsResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'cdn.ListCdnDomains': %w", err)
		}

		if listCdnDomainsResp.Data != nil {
			for _, domain := range listCdnDomainsResp.Data {
				if xcerthostname.IsMatch(wildcardDomain, ve.StringValue(domain.Domain)) {
					domains = append(domains, ve.StringValue(domain.Domain))
				}
			}
		}

		if len(listCdnDomainsResp.Data) < int(listCdnDomainsPageSize) {
			break
		} else {
			listCdnDomainsPageSize++
		}
	}

	return domains, nil
}

func (d *SSLDeployerProvider) getMatchedDomainsByCertId(ctx context.Context, cloudCertId string) ([]string, error) {
	domains := make([]string, 0)

	// 获取指定证书可关联的域名
	// REF: https://www.volcengine.com/docs/6454/125711
	describeCertConfigReq := &vecdn.DescribeCertConfigInput{
		CertId: ve.String(cloudCertId),
	}
	describeCertConfigResp, err := d.sdkClient.DescribeCertConfig(describeCertConfigReq)
	d.logger.Debug("sdk request 'cdn.DescribeCertConfig'", slog.Any("request", describeCertConfigReq), slog.Any("response", describeCertConfigResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'cdn.DescribeCertConfig': %w", err)
	}

	if describeCertConfigResp.CertNotConfig != nil {
		for i := range describeCertConfigResp.CertNotConfig {
			domains = append(domains, ve.StringValue(describeCertConfigResp.CertNotConfig[i].Domain))
		}
	}

	if describeCertConfigResp.OtherCertConfig != nil {
		for i := range describeCertConfigResp.OtherCertConfig {
			domains = append(domains, ve.StringValue(describeCertConfigResp.OtherCertConfig[i].Domain))
		}
	}

	if len(domains) == 0 {
		if len(describeCertConfigResp.SpecifiedCertConfig) == 0 {
			return nil, errors.New("domains not found")
		}
	}

	return domains, nil
}

func (d *SSLDeployerProvider) bindCert(ctx context.Context, domain string, cloudCertId string) error {
	// 关联证书与加速域名
	// REF: https://www.volcengine.com/docs/6454/125712
	batchDeployCertReq := &vecdn.BatchDeployCertInput{
		Domain: ve.String(domain),
		CertId: ve.String(cloudCertId),
	}
	batchDeployCertResp, err := d.sdkClient.BatchDeployCert(batchDeployCertReq)
	d.logger.Debug("sdk request 'cdn.BatchDeployCert'", slog.Any("request", batchDeployCertReq), slog.Any("response", batchDeployCertResp))
	if err != nil {
		return err
	}

	return nil
}

func createSDKClient(accessKeyId, accessKeySecret string) (vecdn.CDNAPI, error) {
	config := ve.NewConfig().
		WithAkSk(accessKeyId, accessKeySecret)

	session, err := vesession.NewSession(config)
	if err != nil {
		return nil, err
	}

	client := vecdn.New(session)
	return client, nil
}
