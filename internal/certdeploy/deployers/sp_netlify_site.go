package deployers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core/deployer"
	netlifysite "github.com/certimate-go/certimate/pkg/core/deployer/providers/netlify-site"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	Registries.MustRegister(domain.DeploymentProviderTypeNetlifySite, func(options *ProviderFactoryOptions) (deployer.Provider, error) {
		credentials := domain.AccessConfigForNetlify{}
		if err := xmaps.Populate(options.ProviderAccessConfig, &credentials); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := netlifysite.NewDeployer(&netlifysite.DeployerConfig{
			ApiToken: credentials.ApiToken,
			SiteId:   xmaps.GetString(options.ProviderExtendedConfig, "siteId"),
		})
		return provider, err
	})
}
