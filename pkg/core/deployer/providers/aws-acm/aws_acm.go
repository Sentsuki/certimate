package awsacm

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	awscred "github.com/aws/aws-sdk-go-v2/credentials"
	awsacm "github.com/aws/aws-sdk-go-v2/service/acm"

	"github.com/certimate-go/certimate/pkg/core/certmgr"
	mcertmgr "github.com/certimate-go/certimate/pkg/core/certmgr/providers/aws-acm"
	"github.com/certimate-go/certimate/pkg/core/deployer"
	xcert "github.com/certimate-go/certimate/pkg/utils/cert"
)

type DeployerConfig struct {
	// AWS AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// AWS SecretAccessKey。
	SecretAccessKey string `json:"secretAccessKey"`
	// AWS 区域。
	Region string `json:"region"`
	// ACM 证书 ARN。
	// 选填。零值时表示新建证书；否则表示更新证书。
	CertificateArn string `json:"certificateArn,omitempty"`
}

type Deployer struct {
	config     *DeployerConfig
	logger     *slog.Logger
	sdkClient  *awsacm.Client
	sdkCertmgr certmgr.Provider
}

var _ deployer.Provider = (*Deployer)(nil)

func NewDeployer(config *DeployerConfig) (*Deployer, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	client, err := createSDKClient(config.AccessKeyId, config.SecretAccessKey, config.Region)
	if err != nil {
		return nil, fmt.Errorf("could not create sdk client: %w", err)
	}

	pcertmgr, err := mcertmgr.NewCertmgr(&mcertmgr.CertmgrConfig{
		AccessKeyId:     config.AccessKeyId,
		SecretAccessKey: config.SecretAccessKey,
		Region:          config.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create ssl manager: %w", err)
	}

	return &Deployer{
		config:     config,
		logger:     slog.Default(),
		sdkClient:  client,
		sdkCertmgr: pcertmgr,
	}, nil
}

func (d *Deployer) SetLogger(logger *slog.Logger) {
	if logger == nil {
		d.logger = slog.New(slog.DiscardHandler)
	} else {
		d.logger = logger
	}

	d.sdkCertmgr.SetLogger(logger)
}

func (d *Deployer) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*deployer.DeployResult, error) {
	if d.config.CertificateArn == "" {
		// 上传证书
		upres, err := d.sdkCertmgr.Upload(ctx, certPEM, privkeyPEM)
		if err != nil {
			return nil, fmt.Errorf("failed to upload certificate file: %w", err)
		} else {
			d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
		}
	} else {
		// 提取服务器证书
		serverCertPEM, intermediaCertPEM, err := xcert.ExtractCertificatesFromPEM(certPEM)
		if err != nil {
			return nil, fmt.Errorf("failed to extract certs: %w", err)
		}

		// 导入证书
		// REF: https://docs.aws.amazon.com/en_us/acm/latest/APIReference/API_ImportCertificate.html
		importCertificateReq := &awsacm.ImportCertificateInput{
			CertificateArn:   aws.String(d.config.CertificateArn),
			Certificate:      ([]byte)(serverCertPEM),
			CertificateChain: ([]byte)(intermediaCertPEM),
			PrivateKey:       ([]byte)(privkeyPEM),
		}
		importCertificateResp, err := d.sdkClient.ImportCertificate(ctx, importCertificateReq)
		d.logger.Debug("sdk request 'acm.ImportCertificate'", slog.Any("request", importCertificateReq), slog.Any("response", importCertificateResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'acm.ImportCertificate': %w", err)
		}
	}

	return &deployer.DeployResult{}, nil
}

func createSDKClient(accessKeyId, secretAccessKey, region string) (*awsacm.Client, error) {
	cfg, err := awscfg.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	client := awsacm.NewFromConfig(cfg, func(o *awsacm.Options) {
		o.Region = region
		o.Credentials = aws.NewCredentialsCache(awscred.NewStaticCredentialsProvider(accessKeyId, secretAccessKey, ""))
	})
	return client, nil
}
