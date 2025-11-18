package tencentcloudeo

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/samber/lo"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tcteo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/certimate-go/certimate/pkg/core/certmgr"
	mcertmgr "github.com/certimate-go/certimate/pkg/core/certmgr/providers/tencentcloud-ssl"
	"github.com/certimate-go/certimate/pkg/core/deployer"
	"github.com/certimate-go/certimate/pkg/core/deployer/providers/tencentcloud-eo/internal"
	xcert "github.com/certimate-go/certimate/pkg/utils/cert"
	xcerthostname "github.com/certimate-go/certimate/pkg/utils/cert/hostname"
)

type DeployerConfig struct {
	// 腾讯云 SecretId。
	SecretId string `json:"secretId"`
	// 腾讯云 SecretKey。
	SecretKey string `json:"secretKey"`
	// 腾讯云接口端点。
	Endpoint string `json:"endpoint,omitempty"`
	// 站点 ID。
	ZoneId string `json:"zoneId"`
	// 域名匹配模式。
	// 零值时默认值 [DOMAIN_MATCH_PATTERN_EXACT]。
	DomainMatchPattern string `json:"domainMatchPattern,omitempty"`
	// 加速域名列表（支持泛域名）。
	Domains []string `json:"domains"`
}

type Deployer struct {
	config     *DeployerConfig
	logger     *slog.Logger
	sdkClient  *internal.TeoClient
	sdkCertmgr certmgr.Provider
}

var _ deployer.Provider = (*Deployer)(nil)

func NewDeployer(config *DeployerConfig) (*Deployer, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	client, err := createSDKClient(config.SecretId, config.SecretKey, config.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("could not create sdk client: %w", err)
	}

	pcertmgr, err := mcertmgr.NewCertmgr(&mcertmgr.CertmgrConfig{
		SecretId:  config.SecretId,
		SecretKey: config.SecretKey,
		Endpoint: lo.
			If(strings.HasSuffix(config.Endpoint, "intl.tencentcloudapi.com"), "ssl.intl.tencentcloudapi.com"). // 国际站使用独立的接口端点
			Else(""),
	})
	if err != nil {
		return nil, fmt.Errorf("could not create ssl manager: %w", err)
	}

	return &Deployer{
		config:     config,
		logger:     slog.Default(),
		sdkClient:  client,
		sdkCertmgr: pcertmgr,
	}, nil
}

func (d *Deployer) SetLogger(logger *slog.Logger) {
	if logger == nil {
		d.logger = slog.New(slog.DiscardHandler)
	} else {
		d.logger = logger
	}

	d.sdkCertmgr.SetLogger(logger)
}

func (d *Deployer) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*deployer.DeployResult, error) {
	if d.config.ZoneId == "" {
		return nil, errors.New("config `zoneId` is required")
	}

	// 上传证书
	upres, err := d.sdkCertmgr.Upload(ctx, certPEM, privkeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to upload certificate file: %w", err)
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	// 获取待部署的域名列表
	var domains []string
	switch d.config.DomainMatchPattern {
	case "", DOMAIN_MATCH_PATTERN_EXACT:
		{
			if len(d.config.Domains) == 0 {
				return nil, errors.New("config `domains` is required")
			}

			domains = d.config.Domains
		}

	case DOMAIN_MATCH_PATTERN_WILDCARD:
		{
			if len(d.config.Domains) == 0 {
				return nil, errors.New("config `domains` is required")
			}

			domainCandidates, err := d.getAllDomainsInZone(ctx, d.config.ZoneId)
			if err != nil {
				return nil, err
			}

			domains = lo.Filter(domainCandidates, func(domain string, _ int) bool {
				for _, configDomain := range d.config.Domains {
					if xcerthostname.IsMatch(configDomain, domain) {
						return true
					}
				}
				return false
			})
			if len(domains) == 0 {
				return nil, errors.New("could not find any domains matched by wildcard")
			}
		}

	case DOMAIN_MATCH_PATTERN_CERTSAN:
		{
			certX509, err := xcert.ParseCertificateFromPEM(certPEM)
			if err != nil {
				return nil, err
			}

			domainCandidates, err := d.getAllDomainsInZone(ctx, d.config.ZoneId)
			if err != nil {
				return nil, err
			}

			domains = lo.Filter(domainCandidates, func(domain string, _ int) bool {
				return certX509.VerifyHostname(domain) == nil
			})
			if len(domains) == 0 {
				return nil, errors.New("could not find any domains matched by certificate")
			}
		}

	default:
		return nil, fmt.Errorf("unsupported domain match pattern: '%s'", d.config.DomainMatchPattern)
	}

	// 配置域名证书
	// REF: https://cloud.tencent.com/document/api/1552/80764
	modifyHostsCertificateReq := tcteo.NewModifyHostsCertificateRequest()
	modifyHostsCertificateReq.ZoneId = common.StringPtr(d.config.ZoneId)
	modifyHostsCertificateReq.Mode = common.StringPtr("sslcert")
	modifyHostsCertificateReq.Hosts = common.StringPtrs(domains)
	modifyHostsCertificateReq.ServerCertInfo = []*tcteo.ServerCertInfo{{CertId: common.StringPtr(upres.CertId)}}
	modifyHostsCertificateResp, err := d.sdkClient.ModifyHostsCertificate(modifyHostsCertificateReq)
	d.logger.Debug("sdk request 'teo.ModifyHostsCertificate'", slog.Any("request", modifyHostsCertificateReq), slog.Any("response", modifyHostsCertificateResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'teo.ModifyHostsCertificate': %w", err)
	}

	return &deployer.DeployResult{}, nil
}

func (d *Deployer) getAllDomainsInZone(ctx context.Context, zoneId string) ([]string, error) {
	var domainsInZone []string

	const pageSize = 200
	for offset := 0; ; offset += pageSize {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// 查询加速域名列表
		// REF: https://cloud.tencent.com/document/api/1552/86336
		describeAccelerationDomainsReq := tcteo.NewDescribeAccelerationDomainsRequest()
		describeAccelerationDomainsReq.Limit = common.Int64Ptr(pageSize)
		describeAccelerationDomainsReq.Offset = common.Int64Ptr(int64(offset))
		describeAccelerationDomainsReq.ZoneId = common.StringPtr(zoneId)
		describeAccelerationDomainsResp, err := d.sdkClient.DescribeAccelerationDomains(describeAccelerationDomainsReq)
		d.logger.Debug("sdk request 'teo.DescribeAccelerationDomains'", slog.Any("request", describeAccelerationDomainsReq), slog.Any("response", describeAccelerationDomainsResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'teo.DescribeAccelerationDomains': %w", err)
		}

		for _, domainItem := range describeAccelerationDomainsResp.Response.AccelerationDomains {
			if domainItem == nil || domainItem.DomainName == nil {
				continue
			}

			domainsInZone = append(domainsInZone, *domainItem.DomainName)
		}

		if len(describeAccelerationDomainsResp.Response.AccelerationDomains) < pageSize {
			break
		}
	}

	return domainsInZone, nil
}

func createSDKClient(secretId, secretKey, endpoint string) (*internal.TeoClient, error) {
	credential := common.NewCredential(secretId, secretKey)

	cpf := profile.NewClientProfile()
	if endpoint != "" {
		cpf.HttpProfile.Endpoint = endpoint
	}

	client, err := internal.NewTeoClient(credential, "", cpf)
	if err != nil {
		return nil, err
	}

	return client, nil
}
