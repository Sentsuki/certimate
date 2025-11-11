package deployers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	rainyunrcdn "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/rainyun-rcdn"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	Registries.MustRegister(domain.DeploymentProviderTypeRainYunRCDN, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		credentials := domain.AccessConfigForRainYun{}
		if err := xmaps.Populate(options.ProviderAccessConfig, &credentials); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := rainyunrcdn.NewSSLDeployerProvider(&rainyunrcdn.SSLDeployerProviderConfig{
			ApiKey:     credentials.ApiKey,
			InstanceId: xmaps.GetInt64(options.ProviderExtendedConfig, "instanceId"),
			Domain:     xmaps.GetString(options.ProviderExtendedConfig, "domain"),
		})
		return provider, err
	})
}
