package dogecloud

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/certimate-go/certimate/pkg/core/certmgr"
	dogesdk "github.com/certimate-go/certimate/pkg/sdk3rd/dogecloud"
)

type CertmgrConfig struct {
	// 多吉云 AccessKey。
	AccessKey string `json:"accessKey"`
	// 多吉云 SecretKey。
	SecretKey string `json:"secretKey"`
}

type Certmgr struct {
	config    *CertmgrConfig
	logger    *slog.Logger
	sdkClient *dogesdk.Client
}

var _ certmgr.Provider = (*Certmgr)(nil)

func NewCertmgr(config *CertmgrConfig) (*Certmgr, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl manager provider is nil")
	}

	client, err := createSDKClient(config.AccessKey, config.SecretKey)
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
	// 生成新证书名（需符合多吉云命名规则）
	certName := fmt.Sprintf("certimate-%d", time.Now().UnixMilli())

	// 上传新证书
	// REF: https://docs.dogecloud.com/cdn/api-cert-upload
	uploadSslCertReq := &dogesdk.UploadCdnCertRequest{
		Note:        certName,
		Certificate: certPEM,
		PrivateKey:  privkeyPEM,
	}
	uploadSslCertResp, err := m.sdkClient.UploadCdnCert(uploadSslCertReq)
	m.logger.Debug("sdk request 'cdn.UploadCdnCert'", slog.Any("request", uploadSslCertReq), slog.Any("response", uploadSslCertResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'cdn.UploadCdnCert': %w", err)
	}

	return &certmgr.UploadResult{
		CertId:   fmt.Sprintf("%d", uploadSslCertResp.Data.Id),
		CertName: certName,
	}, nil
}

func createSDKClient(accessKey, secretKey string) (*dogesdk.Client, error) {
	return dogesdk.NewClient(accessKey, secretKey)
}
