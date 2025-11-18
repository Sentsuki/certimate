package onepanelsite

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/certimate-go/certimate/pkg/core/certmgr"
	mcertmgr "github.com/certimate-go/certimate/pkg/core/certmgr/providers/1panel-ssl"
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
	// 子节点名称。
	// 选填。
	NodeName string `json:"nodeName,omitempty"`
	// 部署资源类型。
	ResourceType string `json:"resourceType"`
	// 网站 ID。
	// 部署资源类型为 [RESOURCE_TYPE_WEBSITE] 时必填。
	WebsiteId int64 `json:"websiteId,omitempty"`
	// 证书 ID。
	// 部署资源类型为 [RESOURCE_TYPE_CERTIFICATE] 时必填。
	CertificateId int64 `json:"certificateId,omitempty"`
}

type Deployer struct {
	config     *DeployerConfig
	logger     *slog.Logger
	sdkClient  any
	sdkCertmgr certmgr.Provider
}

var _ deployer.Provider = (*Deployer)(nil)

func NewDeployer(config *DeployerConfig) (*Deployer, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	client, err := createSDKClient(config.ServerUrl, config.ApiVersion, config.ApiKey, config.AllowInsecureConnections, config.NodeName)
	if err != nil {
		return nil, fmt.Errorf("could not create sdk client: %w", err)
	}

	pcertmgr, err := mcertmgr.NewCertmgr(&mcertmgr.CertmgrConfig{
		ServerUrl:                config.ServerUrl,
		ApiVersion:               config.ApiVersion,
		ApiKey:                   config.ApiKey,
		AllowInsecureConnections: config.AllowInsecureConnections,
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
	// 根据部署资源类型决定部署方式
	switch d.config.ResourceType {
	case RESOURCE_TYPE_WEBSITE:
		if err := d.deployToWebsite(ctx, certPEM, privkeyPEM); err != nil {
			return nil, err
		}

	case RESOURCE_TYPE_CERTIFICATE:
		if err := d.deployToCertificate(ctx, certPEM, privkeyPEM); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported resource type '%s'", d.config.ResourceType)
	}

	return &deployer.DeployResult{}, nil
}

func (d *Deployer) deployToWebsite(ctx context.Context, certPEM string, privkeyPEM string) error {
	if d.config.WebsiteId == 0 {
		return errors.New("config `websiteId` is required")
	}

	// 上传证书
	upres, err := d.sdkCertmgr.Upload(ctx, certPEM, privkeyPEM)
	if err != nil {
		return fmt.Errorf("failed to upload certificate file: %w", err)
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	switch sdkClient := d.sdkClient.(type) {
	case *opsdk.Client:
		{
			// 获取网站 HTTPS 配置
			getHttpsConfResp, err := sdkClient.GetHttpsConf(d.config.WebsiteId)
			d.logger.Debug("sdk request '1panel.GetHttpsConf'", slog.Int64("websiteId", d.config.WebsiteId), slog.Any("response", getHttpsConfResp))
			if err != nil {
				return fmt.Errorf("failed to execute sdk request '1panel.GetHttpsConf': %w", err)
			}

			// 修改网站 HTTPS 配置
			certId, _ := strconv.ParseInt(upres.CertId, 10, 64)
			updateHttpsConfReq := &opsdk.UpdateHttpsConfRequest{
				WebsiteID:    d.config.WebsiteId,
				Type:         "existed",
				WebsiteSSLID: certId,
				Enable:       getHttpsConfResp.Data.Enable,
				HttpConfig:   getHttpsConfResp.Data.HttpConfig,
				SSLProtocol:  getHttpsConfResp.Data.SSLProtocol,
				Algorithm:    getHttpsConfResp.Data.Algorithm,
				Hsts:         getHttpsConfResp.Data.Hsts,
			}
			updateHttpsConfResp, err := sdkClient.UpdateHttpsConf(d.config.WebsiteId, updateHttpsConfReq)
			d.logger.Debug("sdk request '1panel.UpdateHttpsConf'", slog.Int64("websiteId", d.config.WebsiteId), slog.Any("request", updateHttpsConfReq), slog.Any("response", updateHttpsConfResp))
			if err != nil {
				return fmt.Errorf("failed to execute sdk request '1panel.UpdateHttpsConf': %w", err)
			}
		}

	case *opsdkv2.Client:
		{
			// 获取网站 HTTPS 配置
			getHttpsConfResp, err := sdkClient.GetHttpsConf(d.config.WebsiteId)
			d.logger.Debug("sdk request '1panel.GetHttpsConf'", slog.Int64("websiteId", d.config.WebsiteId), slog.Any("response", getHttpsConfResp))
			if err != nil {
				return fmt.Errorf("failed to execute sdk request '1panel.GetHttpsConf': %w", err)
			}

			// 修改网站 HTTPS 配置
			certId, _ := strconv.ParseInt(upres.CertId, 10, 64)
			updateHttpsConfReq := &opsdkv2.UpdateHttpsConfRequest{
				WebsiteID:    d.config.WebsiteId,
				Type:         "existed",
				WebsiteSSLID: certId,
				Enable:       getHttpsConfResp.Data.Enable,
				HttpConfig:   getHttpsConfResp.Data.HttpConfig,
				SSLProtocol:  getHttpsConfResp.Data.SSLProtocol,
				Algorithm:    getHttpsConfResp.Data.Algorithm,
				Hsts:         getHttpsConfResp.Data.Hsts,
			}
			updateHttpsConfResp, err := sdkClient.UpdateHttpsConf(d.config.WebsiteId, updateHttpsConfReq)
			d.logger.Debug("sdk request '1panel.UpdateHttpsConf'", slog.Int64("websiteId", d.config.WebsiteId), slog.Any("request", updateHttpsConfReq), slog.Any("response", updateHttpsConfResp))
			if err != nil {
				return fmt.Errorf("failed to execute sdk request '1panel.UpdateHttpsConf': %w", err)
			}
		}

	default:
		panic("sdk client is not implemented")
	}

	return nil
}

func (d *Deployer) deployToCertificate(ctx context.Context, certPEM string, privkeyPEM string) error {
	if d.config.CertificateId == 0 {
		return errors.New("config `certificateId` is required")
	}

	switch sdkClient := d.sdkClient.(type) {
	case *opsdk.Client:
		{
			// 获取证书详情
			getWebsiteSSLResp, err := sdkClient.GetWebsiteSSL(d.config.CertificateId)
			d.logger.Debug("sdk request '1panel.GetWebsiteSSL'", slog.Int64("sslId", d.config.CertificateId), slog.Any("response", getWebsiteSSLResp))
			if err != nil {
				return fmt.Errorf("failed to execute sdk request '1panel.GetWebsiteSSL': %w", err)
			}

			// 更新证书
			uploadWebsiteSSLReq := &opsdk.UploadWebsiteSSLRequest{
				SSLID:       d.config.CertificateId,
				Type:        "paste",
				Description: getWebsiteSSLResp.Data.Description,
				Certificate: certPEM,
				PrivateKey:  privkeyPEM,
			}
			uploadWebsiteSSLResp, err := sdkClient.UploadWebsiteSSL(uploadWebsiteSSLReq)
			d.logger.Debug("sdk request '1panel.UploadWebsiteSSL'", slog.Any("request", uploadWebsiteSSLReq), slog.Any("response", uploadWebsiteSSLResp))
			if err != nil {
				return fmt.Errorf("failed to execute sdk request '1panel.UploadWebsiteSSL': %w", err)
			}
		}

	case *opsdkv2.Client:
		{
			// 获取证书详情
			getWebsiteSSLResp, err := sdkClient.GetWebsiteSSL(d.config.CertificateId)
			d.logger.Debug("sdk request '1panel.GetWebsiteSSL'", slog.Any("sslId", d.config.CertificateId), slog.Any("response", getWebsiteSSLResp))
			if err != nil {
				return fmt.Errorf("failed to execute sdk request '1panel.GetWebsiteSSL': %w", err)
			}

			// 更新证书
			uploadWebsiteSSLReq := &opsdkv2.UploadWebsiteSSLRequest{
				SSLID:       d.config.CertificateId,
				Type:        "paste",
				Description: getWebsiteSSLResp.Data.Description,
				Certificate: certPEM,
				PrivateKey:  privkeyPEM,
			}
			uploadWebsiteSSLResp, err := sdkClient.UploadWebsiteSSL(uploadWebsiteSSLReq)
			d.logger.Debug("sdk request '1panel.UploadWebsiteSSL'", slog.Any("request", uploadWebsiteSSLReq), slog.Any("response", uploadWebsiteSSLResp))
			if err != nil {
				return fmt.Errorf("failed to execute sdk request '1panel.UploadWebsiteSSL': %w", err)
			}
		}

	default:
		panic("sdk client is not implemented")
	}

	return nil
}

const (
	sdkVersionV1 = "v1"
	sdkVersionV2 = "v2"
)

func createSDKClient(serverUrl, apiVersion, apiKey string, skipTlsVerify bool, nodeName string) (any, error) {
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
		var client *opsdkv2.Client
		if nodeName == "" {
			temp, err := opsdkv2.NewClient(serverUrl, apiKey)
			if err != nil {
				return nil, err
			}
			client = temp
		} else {
			temp, err := opsdkv2.NewClientWithNode(serverUrl, apiKey, nodeName)
			if err != nil {
				return nil, err
			}
			client = temp
		}

		if skipTlsVerify {
			client.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})
		}

		return client, nil
	}

	return nil, fmt.Errorf("invalid 1panel api version")
}
