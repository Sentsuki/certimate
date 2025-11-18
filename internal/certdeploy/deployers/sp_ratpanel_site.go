package deployers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core/deployer"
	ratpanelsite "github.com/certimate-go/certimate/pkg/core/deployer/providers/ratpanel-site"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	Registries.MustRegister(domain.DeploymentProviderTypeRatPanelSite, func(options *ProviderFactoryOptions) (deployer.Provider, error) {
		credentials := domain.AccessConfigForRatPanel{}
		if err := xmaps.Populate(options.ProviderAccessConfig, &credentials); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := ratpanelsite.NewDeployer(&ratpanelsite.DeployerConfig{
			ServerUrl:                credentials.ServerUrl,
			AccessTokenId:            credentials.AccessTokenId,
			AccessToken:              credentials.AccessToken,
			AllowInsecureConnections: credentials.AllowInsecureConnections,
			SiteName:                 xmaps.GetString(options.ProviderExtendedConfig, "siteName"),
		})
		return provider, err
	})
}
