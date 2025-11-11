package applicators

import (
	"github.com/go-acme/lego/v4/challenge"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core/ssl-applicator/acme-http01/providers/local"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

func init() {
	ACMEHttp01Registries.MustRegister(domain.ACMEHttp01ProviderTypeLocal, func(options *ProviderFactoryOptions) (challenge.Provider, error) {
		provider, err := local.NewChallengeProvider(&local.ChallengeProviderConfig{
			WebRootPath: xmaps.GetString(options.ProviderExtendedConfig, "webRootPath"),
		})
		return provider, err
	})
}
