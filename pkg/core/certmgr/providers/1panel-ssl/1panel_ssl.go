package onepanelssl

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/certimate-go/certimate/pkg/core/certmgr"
	opsdk "github.com/certimate-go/certimate/pkg/sdk3rd/1panel"
	opsdkv2 "github.com/certimate-go/certimate/pkg/sdk3rd/1panel/v2"
)

type CertmgrConfig struct {
	// 1Panel 服务地址。
	ServerUrl string `json:"serverUrl"`
	// 1Panel 版本。
	ApiVersion string `json:"apiVersion"`
	// 1Panel 接口密钥。
	ApiKey string `json:"apiKey"`
	// 是否允许不安全的连接。
	AllowInsecureConnections bool `json:"allowInsecureConnections,omitempty"`
}

type Certmgr struct {
	config    *CertmgrConfig
	logger    *slog.Logger
	sdkClient any
}

var _ certmgr.Provider = (*Certmgr)(nil)

func NewCertmgr(config *CertmgrConfig) (*Certmgr, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl manager provider is nil")
	}

	client, err := createSDKClient(config.ServerUrl, config.ApiVersion, config.ApiKey, config.AllowInsecureConnections)
	if err != nil {
		return nil, fmt.Errorf("could not create sdk client: %w", err)
	}

	return &Certmgr{
		config:    config,
		logger:    slog.Default(),
		sdkClient: client,
	}, nil
}

func (m *Certmgr) SetLogger(logger *slog.Logger) {
	if logger == nil {
		m.logger = slog.New(slog.DiscardHandler)
	} else {
		m.logger = logger
	}
}

func (m *Certmgr) Upload(ctx context.Context, certPEM string, privkeyPEM string) (*certmgr.UploadResult, error) {
	// 避免重复上传
	if res, err := m.tryFindCert(ctx, certPEM, privkeyPEM); err != nil {
		return nil, err
	} else if res != nil {
		m.logger.Info("ssl certificate already exists")
		return res, nil
	}

	// 生成新证书名（需符合 1Panel 命名规则）
	certName := fmt.Sprintf("certimate-%d", time.Now().UnixMilli())

	// 上传证书
	switch sdkClient := m.sdkClient.(type) {
	case *opsdk.Client:
		{
			uploadWebsiteSSLReq := &opsdk.UploadWebsiteSSLRequest{
				Type:        "paste",
				Description: certName,
				Certificate: certPEM,
				PrivateKey:  privkeyPEM,
			}
			uploadWebsiteSSLResp, err := sdkClient.UploadWebsiteSSL(uploadWebsiteSSLReq)
			m.logger.Debug("sdk request '1panel.UploadWebsiteSSL'", slog.Any("request", uploadWebsiteSSLReq), slog.Any("response", uploadWebsiteSSLResp))
			if err != nil {
				return nil, fmt.Errorf("failed to execute sdk request '1panel.UploadWebsiteSSL': %w", err)
			}
		}

	case *opsdkv2.Client:
		{
			uploadWebsiteSSLReq := &opsdkv2.UploadWebsiteSSLRequest{
				Type:        "paste",
				Description: certName,
				Certificate: certPEM,
				PrivateKey:  privkeyPEM,
			}
			uploadWebsiteSSLResp, err := sdkClient.UploadWebsiteSSL(uploadWebsiteSSLReq)
			m.logger.Debug("sdk request '1panel.UploadWebsiteSSL'", slog.Any("request", uploadWebsiteSSLReq), slog.Any("response", uploadWebsiteSSLResp))
			if err != nil {
				return nil, fmt.Errorf("failed to execute sdk request '1panel.UploadWebsiteSSL': %w", err)
			}
		}

	default:
		panic("sdk client is not implemented")
	}

	// 获取刚刚上传证书 ID
	if res, err := m.tryFindCert(ctx, certPEM, privkeyPEM); err != nil {
		return nil, err
	} else if res == nil {
		return nil, fmt.Errorf("could not find ssl certificate, may be upload failed")
	} else {
		return res, nil
	}
}

func (m *Certmgr) tryFindCert(ctx context.Context, certPEM string, privkeyPEM string) (*certmgr.UploadResult, error) {
	searchWebsiteSSLPageNumber := 1
	searchWebsiteSSLPageSize := 100
	searchWebsiteSSLItemsCount := 0
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		switch sdkClient := m.sdkClient.(type) {
		case *opsdk.Client:
			{
				searchWebsiteSSLReq := &opsdk.SearchWebsiteSSLRequest{
					Page:     int32(searchWebsiteSSLPageNumber),
					PageSize: int32(searchWebsiteSSLPageSize),
				}
				searchWebsiteSSLResp, err := sdkClient.SearchWebsiteSSL(searchWebsiteSSLReq)
				m.logger.Debug("sdk request '1panel.SearchWebsiteSSL'", slog.Any("request", searchWebsiteSSLReq), slog.Any("response", searchWebsiteSSLResp))
				if err != nil {
					return nil, fmt.Errorf("failed to execute sdk request '1panel.SearchWebsiteSSL': %w", err)
				}

				if searchWebsiteSSLResp.Data == nil {
					break
				}

				for _, sslItem := range searchWebsiteSSLResp.Data.Items {
					if strings.TrimSpace(sslItem.PEM) == strings.TrimSpace(certPEM) &&
						strings.TrimSpace(sslItem.PrivateKey) == strings.TrimSpace(privkeyPEM) {
						// 如果已存在相同证书，直接返回
						return &certmgr.UploadResult{
							CertId:   fmt.Sprintf("%d", sslItem.ID),
							CertName: sslItem.Description,
						}, nil
					}
				}

				searchWebsiteSSLItemsCount = int(searchWebsiteSSLResp.Data.Total)
			}

		case *opsdkv2.Client:
			{
				searchWebsiteSSLReq := &opsdkv2.SearchWebsiteSSLRequest{
					Page:     int32(searchWebsiteSSLPageNumber),
					PageSize: int32(searchWebsiteSSLPageSize),
				}
				searchWebsiteSSLResp, err := sdkClient.SearchWebsiteSSL(searchWebsiteSSLReq)
				m.logger.Debug("sdk request '1panel.SearchWebsiteSSL'", slog.Any("request", searchWebsiteSSLReq), slog.Any("response", searchWebsiteSSLResp))
				if err != nil {
					return nil, fmt.Errorf("failed to execute sdk request '1panel.SearchWebsiteSSL': %w", err)
				}

				if searchWebsiteSSLResp.Data == nil {
					break
				}

				for _, sslItem := range searchWebsiteSSLResp.Data.Items {
					if strings.TrimSpace(sslItem.PEM) == strings.TrimSpace(certPEM) &&
						strings.TrimSpace(sslItem.PrivateKey) == strings.TrimSpace(privkeyPEM) {
						// 如果已存在相同证书，直接返回
						return &certmgr.UploadResult{
							CertId:   fmt.Sprintf("%d", sslItem.ID),
							CertName: sslItem.Description,
						}, nil
					}
				}

				searchWebsiteSSLItemsCount = int(searchWebsiteSSLResp.Data.Total)
			}

		default:
			panic("sdk client is not implemented")
		}

		if searchWebsiteSSLItemsCount < searchWebsiteSSLPageSize {
			break
		} else {
			searchWebsiteSSLPageNumber++
		}
	}

	return nil, nil
}

const (
	sdkVersionV1 = "v1"
	sdkVersionV2 = "v2"
)

func createSDKClient(serverUrl, apiVersion, apiKey string, skipTlsVerify bool) (any, error) {
	if apiVersion == sdkVersionV1 {
		client, err := opsdk.NewClient(serverUrl, apiKey)
		if err != nil {
			return nil, err
		}

		if skipTlsVerify {
			client.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})
		}

		return client, nil
	} else if apiVersion == sdkVersionV2 {
		client, err := opsdkv2.NewClient(serverUrl, apiKey)
		if err != nil {
			return nil, err
		}

		if skipTlsVerify {
			client.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})
		}

		return client, nil
	}

	return nil, fmt.Errorf("invalid 1panel api version")
}
