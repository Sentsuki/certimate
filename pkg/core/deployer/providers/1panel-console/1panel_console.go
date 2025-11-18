package onepanelconsole

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/certimate-go/certimate/pkg/core/deployer"
	opsdk "github.com/certimate-go/certimate/pkg/sdk3rd/1panel"
	opsdkv2 "github.com/certimate-go/certimate/pkg/sdk3rd/1panel/v2"
)

type DeployerConfig struct {
	// 1Panel 服务地址。
	ServerUrl string `json:"serverUrl"`
	// 1Panel 版本。
	// 可取值 "v1"、"v2"。
	ApiVersion string `json:"apiVersion"`
	// 1Panel 接口密钥。
	ApiKey string `json:"apiKey"`
	// 是否允许不安全的连接。
	AllowInsecureConnections bool `json:"allowInsecureConnections,omitempty"`
	// 是否自动重启。
	AutoRestart bool `json:"autoRestart"`
}

type Deployer struct {
	config    *DeployerConfig
	logger    *slog.Logger
	sdkClient any
}

var _ deployer.Provider = (*Deployer)(nil)

func NewDeployer(config *DeployerConfig) (*Deployer, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	client, err := createSDKClient(config.ServerUrl, config.ApiVersion, config.ApiKey, config.AllowInsecureConnections)
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
	// 设置面板 SSL 证书
	switch sdkClient := d.sdkClient.(type) {
	case *opsdk.Client:
		{
			updateSettingsSSLReq := &opsdk.UpdateSettingsSSLRequest{
				Cert:        certPEM,
				Key:         privkeyPEM,
				SSL:         "enable",
				SSLType:     "import-paste",
				AutoRestart: strconv.FormatBool(d.config.AutoRestart),
			}
			updateSystemSSLResp, err := sdkClient.UpdateSettingsSSL(updateSettingsSSLReq)
			d.logger.Debug("sdk request '1panel.UpdateSettingsSSL'", slog.Any("request", updateSettingsSSLReq), slog.Any("response", updateSystemSSLResp))
			if err != nil {
				return nil, fmt.Errorf("failed to execute sdk request '1panel.UpdateSettingsSSL': %w", err)
			}
		}

	case *opsdkv2.Client:
		{
			updateCoreSettingsSSLReq := &opsdkv2.UpdateCoreSettingsSSLRequest{
				Cert:        certPEM,
				Key:         privkeyPEM,
				SSL:         "Enable",
				SSLType:     "import-paste",
				AutoRestart: strconv.FormatBool(d.config.AutoRestart),
			}
			updateCoreSystemSSLResp, err := sdkClient.UpdateCoreSettingsSSL(updateCoreSettingsSSLReq)
			d.logger.Debug("sdk request '1panel.UpdateCoreSettingsSSL'", slog.Any("request", updateCoreSettingsSSLReq), slog.Any("response", updateCoreSystemSSLResp))
			if err != nil {
				return nil, fmt.Errorf("failed to execute sdk request '1panel.UpdateCoreSettingsSSL': %w", err)
			}
		}

	default:
		panic("sdk client is not implemented")
	}

	return &deployer.DeployResult{}, nil
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
