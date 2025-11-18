package upyunssl

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/certimate-go/certimate/pkg/core/certmgr"
	upyunsdk "github.com/certimate-go/certimate/pkg/sdk3rd/upyun/console"
)

type CertmgrConfig struct {
	// 又拍云账号用户名。
	Username string `json:"username"`
	// 又拍云账号密码。
	Password string `json:"password"`
}

type Certmgr struct {
	config    *CertmgrConfig
	logger    *slog.Logger
	sdkClient *upyunsdk.Client
}

var _ certmgr.Provider = (*Certmgr)(nil)

func NewCertmgr(config *CertmgrConfig) (*Certmgr, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl manager provider is nil")
	}

	client, err := createSDKClient(config.Username, config.Password)
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
	// 上传证书
	uploadHttpsCertificateReq := &upyunsdk.UploadHttpsCertificateRequest{
		Certificate: certPEM,
		PrivateKey:  privkeyPEM,
	}
	uploadHttpsCertificateResp, err := m.sdkClient.UploadHttpsCertificate(uploadHttpsCertificateReq)
	m.logger.Debug("sdk request 'console.UploadHttpsCertificate'", slog.Any("response", uploadHttpsCertificateResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'console.UploadHttpsCertificate': %w", err)
	}

	return &certmgr.UploadResult{
		CertId: uploadHttpsCertificateResp.Data.Result.CertificateId,
	}, nil
}

func createSDKClient(username, password string) (*upyunsdk.Client, error) {
	return upyunsdk.NewClient(username, password)
}
