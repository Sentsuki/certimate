package baotapanelwaf

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"

	"github.com/samber/lo"

	"github.com/certimate-go/certimate/pkg/core/deployer"
	btwafsdk "github.com/certimate-go/certimate/pkg/sdk3rd/btwaf"
)

type DeployerConfig struct {
	// 堡塔云 WAF 服务地址。
	ServerUrl string `json:"serverUrl"`
	// 堡塔云 WAF 接口密钥。
	ApiKey string `json:"apiKey"`
	// 是否允许不安全的连接。
	AllowInsecureConnections bool `json:"allowInsecureConnections,omitempty"`
	// 网站名称。
	SiteName string `json:"siteName"`
	// 网站 SSL 端口。
	// 零值时默认值 443。
	SitePort int32 `json:"sitePort,omitempty"`
}

type Deployer struct {
	config    *DeployerConfig
	logger    *slog.Logger
	sdkClient *btwafsdk.Client
}

var _ deployer.Provider = (*Deployer)(nil)

func NewDeployer(config *DeployerConfig) (*Deployer, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	client, err := createSDKClient(config.ServerUrl, config.ApiKey, config.AllowInsecureConnections)
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
	if d.config.SitePort == 0 {
		d.config.SitePort = 443
	}

	// 获取网站 ID
	siteId, err := d.findSiteIdByName(ctx, d.config.SiteName)
	if err != nil {
		return nil, err
	}

	// 修改站点配置
	modifySiteReq := &btwafsdk.ModifySiteRequest{
		SiteId: lo.ToPtr(siteId),
		Type:   lo.ToPtr("openCert"),
		Server: &btwafsdk.SiteServerInfoMod{
			ListenSSLPorts: lo.ToPtr([]string{fmt.Sprintf("%d", d.config.SitePort)}),
			SSL: &btwafsdk.SiteServerSSLInfo{
				IsSSL:      lo.ToPtr(int32(1)),
				FullChain:  lo.ToPtr(certPEM),
				PrivateKey: lo.ToPtr(privkeyPEM),
			},
		},
	}
	modifySiteResp, err := d.sdkClient.ModifySite(modifySiteReq)
	d.logger.Debug("sdk request 'bt.ModifySite'", slog.Any("request", modifySiteReq), slog.Any("response", modifySiteResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'bt.ModifySite': %w", err)
	}

	return &deployer.DeployResult{}, nil
}

func (d *Deployer) findSiteIdByName(ctx context.Context, siteName string) (string, error) {
	// 查询网站列表
	getSiteListPage := 1
	getSiteListPageSize := 100
	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
		}

		getSiteListReq := &btwafsdk.GetSiteListRequest{
			SiteName: lo.ToPtr(d.config.SiteName),
			Page:     lo.ToPtr(int32(getSiteListPage)),
			PageSize: lo.ToPtr(int32(getSiteListPageSize)),
		}
		getSiteListResp, err := d.sdkClient.GetSiteList(getSiteListReq)
		d.logger.Debug("sdk request 'bt.GetSiteList'", slog.Any("request", getSiteListReq), slog.Any("response", getSiteListResp))
		if err != nil {
			return "", fmt.Errorf("failed to execute sdk request 'bt.GetSiteList': %w", err)
		}

		if getSiteListResp.Result == nil {
			break
		}

		for _, siteItem := range getSiteListResp.Result.List {
			if siteItem.SiteName == d.config.SiteName {
				return siteItem.SiteId, nil
			}
		}

		if len(getSiteListResp.Result.List) < getSiteListPageSize {
			break
		}

		getSiteListPage++
	}

	return "", fmt.Errorf("could not find site '%s'", siteName)
}

func createSDKClient(serverUrl, apiKey string, skipTlsVerify bool) (*btwafsdk.Client, error) {
	client, err := btwafsdk.NewClient(serverUrl, apiKey)
	if err != nil {
		return nil, err
	}

	if skipTlsVerify {
		client.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})
	}

	return client, nil
}
