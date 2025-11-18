package baotapanelgoconsole_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "github.com/certimate-go/certimate/pkg/core/deployer/providers/baotapanelgo-console"
)

var (
	fInputCertPath string
	fInputKeyPath  string
	fServerUrl     string
	fApiKey        string
)

func init() {
	argsPrefix := "CERTIMATE_SSLDEPLOYER_BAOTAPANELGOCONSOLE_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fServerUrl, argsPrefix+"SERVERURL", "", "")
	flag.StringVar(&fApiKey, argsPrefix+"APIKEY", "", "")
}

/*
Shell command to run this test:

	go test -v ./baotapanelgo_console_test.go -args \
	--CERTIMATE_SSLDEPLOYER_BAOTAPANELGOCONSOLE_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_SSLDEPLOYER_BAOTAPANELGOCONSOLE_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_SSLDEPLOYER_BAOTAPANELGOCONSOLE_SERVERURL="http://127.0.0.1:8888" \
	--CERTIMATE_SSLDEPLOYER_BAOTAPANELGOCONSOLE_APIKEY="your-api-key"
*/
func TestDeploy(t *testing.T) {
	flag.Parse()

	t.Run("Deploy", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("INPUTCERTPATH: %v", fInputCertPath),
			fmt.Sprintf("INPUTKEYPATH: %v", fInputKeyPath),
			fmt.Sprintf("SERVERURL: %v", fServerUrl),
			fmt.Sprintf("APIKEY: %v", fApiKey),
		}, "\n"))

		provider, err := provider.NewDeployer(&provider.DeployerConfig{
			ServerUrl:                fServerUrl,
			ApiKey:                   fApiKey,
			AllowInsecureConnections: true,
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
