package deployers

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core/deployer"
	baotapanelgosite "github.com/certimate-go/certimate/pkg/core/deployer/providers/baotapanelgo-site"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	Registries.MustRegister(domain.DeploymentProviderTypeBaotaPanelGoSite, func(options *ProviderFactoryOptions) (deployer.Provider, error) {
		credentials := domain.AccessConfigForBaotaPanelGo{}
		if err := xmaps.Populate(options.ProviderAccessConfig, &credentials); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := baotapanelgosite.NewDeployer(&baotapanelgosite.DeployerConfig{
			ServerUrl:                credentials.ServerUrl,
			ApiKey:                   credentials.ApiKey,
			AllowInsecureConnections: credentials.AllowInsecureConnections,
			SiteType:                 xmaps.GetString(options.ProviderExtendedConfig, "siteType"),
			SiteNames:                lo.Filter(strings.Split(xmaps.GetString(options.ProviderExtendedConfig, "siteNames"), ";"), func(s string, _ int) bool { return s != "" }),
		})
		return provider, err
	})
}
