package deployers

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	ksyuncdn "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/ksyun-cdn"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	Registries.MustRegister(domain.DeploymentProviderTypeKsyunCDN, func(options *ProviderFactoryOptions) (core.SSLDeployer, error) {
		credentials := domain.AccessConfigForKsyun{}
		if err := xmaps.Populate(options.ProviderAccessConfig, &credentials); err != nil {
			return nil, fmt.Errorf("failed to populate provider access config: %w", err)
		}

		provider, err := ksyuncdn.NewSSLDeployerProvider(&ksyuncdn.SSLDeployerProviderConfig{
			AccessKeyId:     credentials.AccessKeyId,
			SecretAccessKey: credentials.SecretAccessKey,
			Domain:          xmaps.GetString(options.ProviderExtendedConfig, "domain"),
			CertificateId:   xmaps.GetString(options.ProviderExtendedConfig, "certificateId"),
		})
		return provider, err
	})
}
