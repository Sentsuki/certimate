package ucloududnr

import (
	"errors"
	"time"

	"github.com/certimate-go/certimate/pkg/core"
	"github.com/certimate-go/certimate/pkg/core/ssl-applicator/acme-dns01/providers/ucloud-udnr/internal"
)

type ChallengeProviderConfig struct {
	PrivateKey            string `json:"privateKey"`
	PublicKey             string `json:"publicKey"`
	ProjectId             string `json:"projectId,omitempty"`
	DnsPropagationTimeout int    `json:"dnsPropagationTimeout,omitempty"`
	DnsTTL                int    `json:"dnsTTL,omitempty"`
}

func NewChallengeProvider(config *ChallengeProviderConfig) (core.ACMEChallenger, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	providerConfig := internal.NewDefaultConfig()
	providerConfig.PrivateKey = config.PrivateKey
	providerConfig.PublicKey = config.PublicKey
	providerConfig.ProjectId = config.ProjectId
	if config.DnsTTL != 0 {
		providerConfig.TTL = config.DnsTTL
	}
	if config.DnsPropagationTimeout != 0 {
		providerConfig.PropagationTimeout = time.Duration(config.DnsPropagationTimeout) * time.Second
	}

	provider, err := internal.NewDNSProviderConfig(providerConfig)
	if err != nil {
		return nil, err
	}

	return provider, nil
}
