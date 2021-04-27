package action

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/redhat-cop/oc-helm/pkg/client"
	"github.com/redhat-cop/oc-helm/pkg/options"
)

type Action interface {
	Setup() error
}

type baseAction struct {
	commandLineOptions *options.CommandLineOption
	helmChartClient    *client.HelmChartClient
	openshiftClient    *client.OpenShiftClient
}

func (b *baseAction) Setup() error {

	if err := b.commandLineOptions.Process(); err != nil {
		return err
	}

	// Attempt to disover console hostname. Otherwise fallback to console-openshift-console.apps
	if b.commandLineOptions.ConsoleHostname == "" {

		openShiftClient, err := client.NewOpenShiftClient(b.commandLineOptions)

		if err != nil {
			return err
		}

		consoleURLConfigMap, err := openShiftClient.DiscoverConsoleURL()

		if err != nil {
			return err
		}

		consoleURL, err := extractConsoleURLFromConfigMap(consoleURLConfigMap)

		if err == nil {
			b.commandLineOptions.ConsoleHostname = consoleURL.Hostname()
		} else {

			serverURL, err := url.Parse(b.commandLineOptions.Server)

			if err != nil {
				return err
			}

			b.commandLineOptions.ConsoleHostname = strings.Replace(serverURL.Hostname(), "api", "console-openshift-console.apps", 1)
		}

	}

	helmChartClient, err := client.NewHelmChartClient(b.commandLineOptions)

	if err != nil {
		return err
	}

	b.helmChartClient = helmChartClient

	return nil
}

func extractConsoleURLFromConfigMap(configMap map[string]interface{}) (*url.URL, error) {

	var consoleURL string

	if configMap != nil {
		data := configMap["data"].(map[string]interface{})

		if data["consoleURL"] != nil {
			consoleURL = data["consoleURL"].(string)
		}
	}

	if consoleURL == "" {
		return nil, fmt.Errorf("Unable to determine Console URL from ConfigMap")
	}

	parsedConsoleURL, err := url.Parse(consoleURL)

	if err != nil {
		return nil, err
	}

	return parsedConsoleURL, nil
}
