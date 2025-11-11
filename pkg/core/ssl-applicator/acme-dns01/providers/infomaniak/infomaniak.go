package infomaniak

import (
	"errors"
	"time"

	"github.com/go-acme/lego/v4/providers/dns/infomaniak"

	"github.com/certimate-go/certimate/pkg/core"
)

type ChallengeProviderConfig struct {
	AccessToken           string `json:"accessToken"`
	DnsPropagationTimeout int    `json:"dnsPropagationTimeout,omitempty"`
	DnsTTL                int    `json:"dnsTTL,omitempty"`
}

func NewChallengeProvider(config *ChallengeProviderConfig) (core.ACMEChallenger, error) {
	if config == nil {
		return nil, errors.New("the configuration of the acme challenge provider is nil")
	}

	providerConfig := infomaniak.NewDefaultConfig()
	providerConfig.AccessToken = config.AccessToken
	if config.DnsPropagationTimeout != 0 {
		providerConfig.PropagationTimeout = time.Duration(config.DnsPropagationTimeout) * time.Second
	}
	if config.DnsTTL != 0 {
		providerConfig.TTL = config.DnsTTL
	}

	provider, err := infomaniak.NewDNSProviderConfig(providerConfig)
	if err != nil {
		return nil, err
	}

	return provider, nil
}
