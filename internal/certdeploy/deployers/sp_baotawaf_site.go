package deployers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core/deployer"
	baotawafsite "github.com/certimate-go/certimate/pkg/core/deployer/providers/baotawaf-site"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	Registries.MustRegister(domain.DeploymentProviderTypeBaotaWAFSite, func(options *ProviderFactoryOptions) (deployer.Provider, error) {
		credentials := domain.AccessConfigForBaotaWAF{}
		if err := xmaps.Populate(options.ProviderAccessConfig, &credentials); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := baotawafsite.NewDeployer(&baotawafsite.DeployerConfig{
			ServerUrl:                credentials.ServerUrl,
			ApiKey:                   credentials.ApiKey,
			AllowInsecureConnections: credentials.AllowInsecureConnections,
			SiteName:                 xmaps.GetString(options.ProviderExtendedConfig, "siteName"),
			SitePort:                 xmaps.GetOrDefaultInt32(options.ProviderExtendedConfig, "sitePort", 443),
		})
		return provider, err
	})
}
