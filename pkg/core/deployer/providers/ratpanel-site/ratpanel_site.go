package ratpanelsite

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"

	"github.com/certimate-go/certimate/pkg/core/deployer"
	rpsdk "github.com/certimate-go/certimate/pkg/sdk3rd/ratpanel"
)

type DeployerConfig struct {
	// 耗子面板服务地址。
	ServerUrl string `json:"serverUrl"`
	// 耗子面板访问令牌 ID。
	AccessTokenId int64 `json:"accessTokenId"`
	// 耗子面板访问令牌。
	AccessToken string `json:"accessToken"`
	// 是否允许不安全的连接。
	AllowInsecureConnections bool `json:"allowInsecureConnections,omitempty"`
	// 网站名称。
	SiteName string `json:"siteName"`
}

type Deployer struct {
	config    *DeployerConfig
	logger    *slog.Logger
	sdkClient *rpsdk.Client
}

var _ deployer.Provider = (*Deployer)(nil)

func NewDeployer(config *DeployerConfig) (*Deployer, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	client, err := createSDKClient(config.ServerUrl, config.AccessTokenId, config.AccessToken, config.AllowInsecureConnections)
	if err != nil {
		return nil, fmt.Errorf("could not create sdk client: %w", err)
	}

	return &Deployer{
		config:    config,
		logger:    slog.Default(),
		sdkClient: client,
	}, nil
}

func (d *Deployer) SetLogger(logger *slog.Logger) {
	if logger == nil {
		d.logger = slog.New(slog.DiscardHandler)
	} else {
		d.logger = logger
	}
}

func (d *Deployer) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*deployer.DeployResult, error) {
	if d.config.SiteName == "" {
		return nil, errors.New("config `siteName` is required")
	}

	// 设置站点 SSL 证书
	setWebsiteCertReq := &rpsdk.SetWebsiteCertRequest{
		SiteName:    d.config.SiteName,
		Certificate: certPEM,
		PrivateKey:  privkeyPEM,
	}
	setWebsiteCertResp, err := d.sdkClient.SetWebsiteCert(setWebsiteCertReq)
	d.logger.Debug("sdk request 'ratpanel.SetWebsiteCert'", slog.Any("request", setWebsiteCertReq), slog.Any("response", setWebsiteCertResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'ratpanel.SetWebsiteCert': %w", err)
	}

	return &deployer.DeployResult{}, nil
}

func createSDKClient(serverUrl string, accessTokenId int64, accessToken string, skipTlsVerify bool) (*rpsdk.Client, error) {
	client, err := rpsdk.NewClient(serverUrl, accessTokenId, accessToken)
	if err != nil {
		return nil, err
	}

	if skipTlsVerify {
		client.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})
	}

	return client, nil
}
