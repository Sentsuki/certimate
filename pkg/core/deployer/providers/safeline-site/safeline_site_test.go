package safelinesite_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "github.com/certimate-go/certimate/pkg/core/deployer/providers/safeline-site"
)

var (
	fInputCertPath string
	fInputKeyPath  string
	fServerUrl     string
	fApiToken      string
	fCertificateId int64
)

func init() {
	argsPrefix := "SAFELINESITE_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fServerUrl, argsPrefix+"SERVERURL", "", "")
	flag.StringVar(&fApiToken, argsPrefix+"APITOKEN", "", "")
	flag.Int64Var(&fCertificateId, argsPrefix+"CERTIFICATEID", 0, "")
}

/*
Shell command to run this test:

	go test -v ./safeline_site_test.go -args \
	--SAFELINESITE_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--SAFELINESITE_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--SAFELINESITE_SERVERURL="http://127.0.0.1:9443" \
	--SAFELINESITE_APITOKEN="your-api-token" \
	--SAFELINESITE_CERTIFICATEID="your-certificate-id"
*/
func TestDeploy(t *testing.T) {
	flag.Parse()

	t.Run("Deploy", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("INPUTCERTPATH: %v", fInputCertPath),
			fmt.Sprintf("INPUTKEYPATH: %v", fInputKeyPath),
			fmt.Sprintf("SERVERURL: %v", fServerUrl),
			fmt.Sprintf("APITOKEN: %v", fApiToken),
			fmt.Sprintf("CERTIFICATEID: %v", fCertificateId),
		}, "\n"))

		provider, err := provider.NewDeployer(&provider.DeployerConfig{
			ServerUrl:                fServerUrl,
			ApiToken:                 fApiToken,
			AllowInsecureConnections: true,
			ResourceType:             provider.RESOURCE_TYPE_CERTIFICATE,
			CertificateId:            fCertificateId,
		})
		if err != nil {
			t.Errorf("err: %+v", err)
			return
		}

		fInputCertData, _ := os.ReadFile(fInputCertPath)
		fInputKeyData, _ := os.ReadFile(fInputKeyPath)
		res, err := provider.Deploy(context.Background(), string(fInputCertData), string(fInputKeyData))
		if err != nil {
			t.Errorf("err: %+v", err)
			return
		}

		t.Logf("ok: %v", res)
	})
}
