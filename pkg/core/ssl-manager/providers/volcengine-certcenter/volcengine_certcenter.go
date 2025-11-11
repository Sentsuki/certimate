package volcenginecertcenter

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	vecs "github.com/volcengine/volcengine-go-sdk/service/certificateservice"
	ve "github.com/volcengine/volcengine-go-sdk/volcengine"
	vesession "github.com/volcengine/volcengine-go-sdk/volcengine/session"

	"github.com/certimate-go/certimate/pkg/core"
)

type SSLManagerProviderConfig struct {
	// 火山引擎 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 火山引擎 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 火山引擎地域。
	Region string `json:"region"`
}

type SSLManagerProvider struct {
	config    *SSLManagerProviderConfig
	logger    *slog.Logger
	sdkClient vecs.CERTIFICATESERVICEAPI
}

var _ core.SSLManager = (*SSLManagerProvider)(nil)

func NewSSLManagerProvider(config *SSLManagerProviderConfig) (*SSLManagerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl manager provider is nil")
	}

	client, err := createSDKClient(config.AccessKeyId, config.AccessKeySecret, config.Region)
	if err != nil {
		return nil, fmt.Errorf("could not create sdk client: %w", err)
	}

	return &SSLManagerProvider{
		config:    config,
		logger:    slog.Default(),
		sdkClient: client,
	}, nil
}

func (m *SSLManagerProvider) SetLogger(logger *slog.Logger) {
	if logger == nil {
		m.logger = slog.New(slog.DiscardHandler)
	} else {
		m.logger = logger
	}
}

func (m *SSLManagerProvider) Upload(ctx context.Context, certPEM string, privkeyPEM string) (*core.SSLManageUploadResult, error) {
	// 上传证书
	// REF: https://www.volcengine.com/docs/6638/1365580
	importCertificateReq := &vecs.ImportCertificateInput{
		CertificateInfo: &vecs.CertificateInfoForImportCertificateInput{
			CertificateChain: ve.String(certPEM),
			PrivateKey:       ve.String(privkeyPEM),
		},
		Repeatable: ve.Bool(false),
	}
	importCertificateResp, err := m.sdkClient.ImportCertificate(importCertificateReq)
	m.logger.Debug("sdk request 'certcenter.ImportCertificate'", slog.Any("request", importCertificateReq), slog.Any("response", importCertificateResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'certcenter.ImportCertificate': %w", err)
	}

	var sslId string
	if importCertificateResp.InstanceId != nil && *importCertificateResp.InstanceId != "" {
		sslId = *importCertificateResp.InstanceId
	}
	if importCertificateResp.RepeatId != nil && *importCertificateResp.RepeatId != "" {
		sslId = *importCertificateResp.RepeatId
	}

	if sslId == "" {
		return nil, errors.New("received empty certificate id, both `InstanceId` and `RepeatId` are empty")
	}

	return &core.SSLManageUploadResult{
		CertId: sslId,
	}, nil
}

func createSDKClient(accessKeyId, accessKeySecret, region string) (vecs.CERTIFICATESERVICEAPI, error) {
	if region == "" {
		region = "cn-beijing" // 证书中心默认区域：北京
	}

	config := ve.NewConfig().
		WithAkSk(accessKeyId, accessKeySecret).
		WithRegion(region)

	session, err := vesession.NewSession(config)
	if err != nil {
		return nil, err
	}

	client := vecs.New(session)
	return client, nil
}
