package ssh

import (
	"github.com/certimate-go/certimate/internal/domain"
)

const (
	AUTH_METHOD_NONE     = "none"
	AUTH_METHOD_PASSWORD = "password"
	AUTH_METHOD_KEY      = "key"
)

const (
	OUTPUT_FORMAT_PEM = string(domain.CertificateFormatTypePEM)
	OUTPUT_FORMAT_PFX = string(domain.CertificateFormatTypePFX)
	OUTPUT_FORMAT_JKS = string(domain.CertificateFormatTypeJKS)
)

