package deployers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	opsite "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/1panel-site"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	Registries.MustRegister(domain.DeploymentProviderType1PanelSite, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		credentials := domain.AccessConfigFor1Panel{}
		if err := xmaps.Populate(options.ProviderAccessConfig, &credentials); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := opsite.NewSSLDeployerProvider(&opsite.SSLDeployerProviderConfig{
			ServerUrl:                credentials.ServerUrl,
			ApiVersion:               credentials.ApiVersion,
			ApiKey:                   credentials.ApiKey,
			AllowInsecureConnections: credentials.AllowInsecureConnections,
			NodeName:                 xmaps.GetString(options.ProviderExtendedConfig, "nodeName"),
			ResourceType:             xmaps.GetString(options.ProviderExtendedConfig, "resourceType"),
			WebsiteId:                xmaps.GetInt64(options.ProviderExtendedConfig, "websiteId"),
			CertificateId:            xmaps.GetInt64(options.ProviderExtendedConfig, "certificateId"),
		})
		return provider, err
	})
}
