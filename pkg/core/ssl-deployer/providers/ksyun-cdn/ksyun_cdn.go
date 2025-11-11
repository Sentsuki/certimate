package ksyuncdn

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/KscSDK/ksc-sdk-go/ksc"
	ksccdnv1 "github.com/KscSDK/ksc-sdk-go/service/cdnv1"
	"github.com/go-viper/mapstructure/v2"

	"github.com/certimate-go/certimate/pkg/core"
)

type SSLDeployerProviderConfig struct {
	// 金山云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 金山云 SecretAccessKey。
	SecretAccessKey string `json:"secretAccessKey"`
	// 加速域名（支持泛域名）。
	Domain string `json:"domain"`
	// 证书 ID。
	// 选填。零值时表示新建证书；否则表示更新证书。
	CertificateId string `json:"certificateId,omitempty"`
}

type SSLDeployerProvider struct {
	config    *SSLDeployerProviderConfig
	logger    *slog.Logger
	sdkClient *ksccdnv1.Cdnv1
}

var _ core.SSLDeployer = (*SSLDeployerProvider)(nil)

func NewSSLDeployerProvider(config *SSLDeployerProviderConfig) (*SSLDeployerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	client, err := createSDKClient(config.AccessKeyId, config.SecretAccessKey)
	if err != nil {
		return nil, fmt.Errorf("could not create sdk client: %w", err)
	}

	return &SSLDeployerProvider{
		config:    config,
		logger:    slog.Default(),
		sdkClient: client,
	}, nil
}

func (d *SSLDeployerProvider) SetLogger(logger *slog.Logger) {
	if logger == nil {
		d.logger = slog.New(slog.DiscardHandler)
	} else {
		d.logger = logger
	}
}

func (d *SSLDeployerProvider) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*core.SSLDeployResult, error) {
	// // 如果原证书 ID 为空，则创建证书；否则更新证书。
	if d.config.CertificateId == "" {
		if d.config.Domain == "" {
			return nil, errors.New("config `domain` is required")
		}

		// 遍历查询域名列表，获取域名 ID
		// https://docs.ksyun.com/documents/198
		var domainId string
		getCdnDomainsPageNumber := int32(1)
		getCdnDomainsPageSize := int32(100)
		for {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}

			getCdnDomainsInput := map[string]any{
				"PageNumber": getCdnDomainsPageNumber,
				"PageSize":   getCdnDomainsPageSize,
				"DomainName": d.config.Domain,
				"FuzzyMatch": "off",
			}
			getCdnDomainsReq, getCdnDomainsOutput := d.sdkClient.GetCdnDomainsPostRequest(&getCdnDomainsInput)
			getCdnDomainsErr := getCdnDomainsReq.Send()
			d.logger.Debug("sdk request 'cdn.GetCdnDomains'", slog.Any("request", getCdnDomainsInput), slog.Any("response", getCdnDomainsOutput))
			if getCdnDomainsErr != nil {
				return nil, fmt.Errorf("failed to execute sdk request 'cdn.GetCdnDomains': %w", getCdnDomainsErr)
			}

			type GetCdnDomainsResponse struct {
				PageNumber int32 `json:"PageNumber"`
				PageSize   int32 `json:"PageSize"`
				TotalCount int32 `json:"TotalCount"`
				Domains    []*struct {
					DomainId     string `json:"DomainId"`
					DomainName   string `json:"DomainName"`
					Cname        string `json:"Cname"`
					CdnType      string `json:"CdnType"`
					CreatedTime  string `json:"CreatedTime"`
					ModifiedTime string `json:"ModifiedTime"`
					Region       string `json:"Region"`
				} `json:"Domains"`
			}
			var getCdnDomainsResp *GetCdnDomainsResponse
			mapstructure.Decode(getCdnDomainsOutput, &getCdnDomainsResp)

			if getCdnDomainsResp != nil {
				for _, domainItem := range getCdnDomainsResp.Domains {
					if strings.EqualFold(domainItem.DomainName, d.config.Domain) {
						domainId = domainItem.DomainId
						break
					}
				}

				if domainId != "" {
					break
				}
			}

			if getCdnDomainsResp == nil || len(getCdnDomainsResp.Domains) < int(getCdnDomainsPageSize) {
				break
			} else {
				getCdnDomainsPageNumber++
			}
		}
		if domainId == "" {
			return nil, errors.New("domain not found")
		}

		// 为加速域名配置证书接口
		// https://docs.ksyun.com/documents/261
		configCertificateInput := map[string]any{
			"Enable":            "on",
			"DomainIds":         domainId,
			"CertificateName":   fmt.Sprintf("certimate_%d", time.Now().UnixMilli()),
			"ServerCertificate": certPEM,
			"PrivateKey":        privkeyPEM,
		}
		configCertificateReq, configCertificateOutput := d.sdkClient.ConfigCertificatePostRequest(&configCertificateInput)
		configCertificateErr := configCertificateReq.Send()
		d.logger.Debug("sdk request 'cdn.ConfigCertificate'", slog.Any("request", configCertificateInput), slog.Any("response", configCertificateOutput))
		if configCertificateErr != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'cdn.ConfigCertificate': %w", configCertificateErr)
		}
	} else {
		// 更新证书
		// https://docs.ksyun.com/documents/259
		setCertificateInput := map[string]any{
			"CertificateId":     d.config.CertificateId,
			"CertificateName":   fmt.Sprintf("certimate_%d", time.Now().UnixMilli()),
			"ServerCertificate": certPEM,
			"PrivateKey":        privkeyPEM,
		}
		setCertificateReq, setCertificateOutput := d.sdkClient.SetCertificatePostRequest(&setCertificateInput)
		setCertificateErr := setCertificateReq.Send()
		d.logger.Debug("sdk request 'cdn.SetCertificate'", slog.Any("request", setCertificateInput), slog.Any("response", setCertificateOutput))
		if setCertificateErr != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'cdn.SetCertificate': %w", setCertificateErr)
		}
	}

	return &core.SSLDeployResult{}, nil
}

func createSDKClient(accessKeyId, secretAccessKey string) (*ksccdnv1.Cdnv1, error) {
	region := "cn-beijing-6"
	client := ksccdnv1.SdkNew(ksc.NewClient(accessKeyId, secretAccessKey), &ksc.Config{Region: &region})
	return client, nil
}
