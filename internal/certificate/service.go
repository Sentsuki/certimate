package certificate

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/x509"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/pocketbase/dbx"

	"github.com/certimate-go/certimate/internal/app"
	"github.com/certimate-go/certimate/internal/certapply"
	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/internal/domain/dtos"
	xcert "github.com/certimate-go/certimate/pkg/utils/cert"
	xcryptokey "github.com/certimate-go/certimate/pkg/utils/crypto/key"
)

type CertificateService struct {
	acmeAccountRepo acmeAccountRepository
	certificateRepo certificateRepository
	settingsRepo    settingsRepository
}

func NewCertificateService(
	acmeAccountRepo acmeAccountRepository,
	certificateRepo certificateRepository,
	settingsRepo settingsRepository,
) *CertificateService {
	return &CertificateService{
		acmeAccountRepo: acmeAccountRepo,
		certificateRepo: certificateRepo,
		settingsRepo:    settingsRepo,
	}
}

func (s *CertificateService) InitSchedule(ctx context.Context) error {
	// 每日清理过期证书
	app.GetScheduler().MustAdd("cleanupCertificateExpired", "0 0 * * *", func() {
		s.cleanupExpiredCertificates(context.Background())
	})

	return nil
}

func (s *CertificateService) DownloadArchivedFile(ctx context.Context, req *dtos.CertificateArchiveFileReq) (*dtos.CertificateArchiveFileResp, error) {
	certificate, err := s.certificateRepo.GetById(ctx, req.CertificateId)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)
	defer zipWriter.Close()

	resp := &dtos.CertificateArchiveFileResp{
		FileFormat: "zip",
	}

	switch strings.ToUpper(req.Format) {
	case "", "PEM":
		{
			certWriter, err := zipWriter.Create("certbundle.pem")
			if err != nil {
				return nil, err
			}

			_, err = certWriter.Write([]byte(certificate.Certificate))
			if err != nil {
				return nil, err
			}

			keyWriter, err := zipWriter.Create("privkey.pem")
			if err != nil {
				return nil, err
			}

			_, err = keyWriter.Write([]byte(certificate.PrivateKey))
			if err != nil {
				return nil, err
			}

			err = zipWriter.Close()
			if err != nil {
				return nil, err
			}

			resp.FileBytes = buf.Bytes()
			return resp, nil
		}

	case "PFX":
		{
			const pfxPassword = "certimate"

			certPFX, err := xcert.TransformCertificateFromPEMToPFX(certificate.Certificate, certificate.PrivateKey, pfxPassword)
			if err != nil {
				return nil, err
			}

			certWriter, err := zipWriter.Create("cert.pfx")
			if err != nil {
				return nil, err
			}

			_, err = certWriter.Write(certPFX)
			if err != nil {
				return nil, err
			}

			keyWriter, err := zipWriter.Create("pfx-password.txt")
			if err != nil {
				return nil, err
			}

			_, err = keyWriter.Write([]byte(pfxPassword))
			if err != nil {
				return nil, err
			}

			err = zipWriter.Close()
			if err != nil {
				return nil, err
			}

			resp.FileBytes = buf.Bytes()
			return resp, nil
		}

	case "JKS":
		{
			const jksPassword = "certimate"

			certJKS, err := xcert.TransformCertificateFromPEMToJKS(certificate.Certificate, certificate.PrivateKey, jksPassword, jksPassword, jksPassword)
			if err != nil {
				return nil, err
			}

			certWriter, err := zipWriter.Create("cert.jks")
			if err != nil {
				return nil, err
			}

			_, err = certWriter.Write(certJKS)
			if err != nil {
				return nil, err
			}

			keyWriter, err := zipWriter.Create("jks-password.txt")
			if err != nil {
				return nil, err
			}

			_, err = keyWriter.Write([]byte(jksPassword))
			if err != nil {
				return nil, err
			}

			err = zipWriter.Close()
			if err != nil {
				return nil, err
			}

			resp.FileBytes = buf.Bytes()
			return resp, nil
		}

	default:
		return nil, domain.ErrInvalidParams
	}
}

func (s *CertificateService) RevokeCertificate(ctx context.Context, req *dtos.CertificateRevokeReq) (*dtos.CertificateRevokeResp, error) {
	certificate, err := s.certificateRepo.GetById(ctx, req.CertificateId)
	if err != nil {
		return nil, err
	}

	if certificate.ACMEAcctUrl == "" || certificate.ACMECertUrl == "" {
		return nil, fmt.Errorf("could not revoke a certificate which is not issued in Certimate")
	}
	// if certificate.ValidityNotAfter.Before(time.Now()) {
	// 	return nil, fmt.Errorf("could not revoke a certificate which is expired")
	// }
	if certificate.IsRevoked {
		return nil, fmt.Errorf("could not revoke a certificate which is already revoked")
	}

	acmeAccount, err := s.acmeAccountRepo.GetByAcctUrl(ctx, certificate.ACMEAcctUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to revoke certificate: could not find acme account: %w", err)
	}

	legoClient, err := certapply.NewACMEClientWithAccount(acmeAccount)
	if err != nil {
		return nil, fmt.Errorf("failed to revoke certificate: could not initialize acme config: %w", err)
	}

	revokeReq := &certapply.RevokeCertificateRequest{
		Certificate: certificate.Certificate,
	}
	_, err = legoClient.RevokeCertificate(ctx, revokeReq)
	if err != nil {
		return nil, fmt.Errorf("failed to revoke certificate: %w", err)
	}

	certificate.IsRevoked = true
	certificate, err = s.certificateRepo.Save(ctx, certificate)
	if err != nil {
		return nil, err
	}

	return &dtos.CertificateRevokeResp{}, nil
}

func (s *CertificateService) ValidateCertificate(ctx context.Context, req *dtos.CertificateValidateCertificateReq) (*dtos.CertificateValidateCertificateResp, error) {
	certX509, err := xcert.ParseCertificateFromPEM(req.Certificate)
	if err != nil {
		return nil, err
	} else if certX509.NotAfter.Before(time.Now()) {
		return nil, fmt.Errorf("certificate has expired at %s", certX509.NotAfter.UTC().Format(time.RFC3339))
	}

	return &dtos.CertificateValidateCertificateResp{
		IsValid: true,
		Domains: strings.Join(certX509.DNSNames, ";"),
	}, nil
}

func (s *CertificateService) ValidatePrivateKey(ctx context.Context, req *dtos.CertificateValidatePrivateKeyReq) (*dtos.CertificateValidatePrivateKeyResp, error) {
	privkey, err := xcert.ParsePrivateKeyFromPEM(req.PrivateKey)
	if err != nil {
		return nil, err
	}

	var keyAlgorithmString string
	keyAlgorithm, keySize, _ := xcryptokey.GetPrivateKeyAlgorithm(privkey)
	switch keyAlgorithm {
	case x509.RSA:
		keyAlgorithmString = fmt.Sprintf("RSA%d", keySize)
	case x509.ECDSA:
		keyAlgorithmString = fmt.Sprintf("EC%d", keySize)
	case x509.Ed25519:
		keyAlgorithmString = "ED25519"
	}

	return &dtos.CertificateValidatePrivateKeyResp{
		IsValid:      keyAlgorithmString != "",
		KeyAlgorithm: keyAlgorithmString,
	}, nil
}

func (s *CertificateService) cleanupExpiredCertificates(ctx context.Context) error {
	settings, err := s.settingsRepo.GetByName(ctx, "persistence")
	if err != nil {
		app.GetLogger().Error("failed to get persistence settings", slog.Any("error", err))
		return err
	}

	persistenceSettings := settings.Content.AsPersistence()
	if persistenceSettings.ExpiredCertificatesMaxDaysRetention != 0 {
		ret, err := s.certificateRepo.DeleteWhere(
			context.Background(),
			dbx.NewExp(fmt.Sprintf("validityNotAfter<DATETIME('now', '-%d days')", persistenceSettings.ExpiredCertificatesMaxDaysRetention)),
		)
		if err != nil {
			app.GetLogger().Error("failed to delete expired certificates", slog.Any("error", err))
			return err
		}

		if ret > 0 {
			app.GetLogger().Info(fmt.Sprintf("cleanup %d expired certificates", ret))
		}
	}

	return nil
}
