package action

import (
	"fmt"
	"strings"

	"github.com/redhat-cop/oc-helm/pkg/options"
)

type UpgradeAsyncAction struct {
	baseAction
}

func NewUpgradeAsyncAction(commandLineOptions *options.CommandLineOption) *UpgradeAsyncAction {
	return &UpgradeAsyncAction{
		baseAction: baseAction{
			commandLineOptions: commandLineOptions,
		},
	}
}

func (u *UpgradeAsyncAction) Run(releaseName string, chartReference string) error {

	_, err := u.helmChartClient.GetRelease(releaseName)

	helmChartInstall := &HelmChartInstall{
		releaseName:        releaseName,
		chartReference:     chartReference,
		commandLineOptions: u.commandLineOptions,
		helmChartClient:    u.helmChartClient,
		upgrade:            true,
	}

	// TODO: Change logic to return HelmClientError to Actions level to inspect status code response
	if err != nil && u.commandLineOptions.Install && strings.Contains(err.Error(), "release: not found") {

		fmt.Fprintf(u.commandLineOptions.Streams.Out, "Release \"%s\" does not exist. Installing it now.\n", releaseName)

		helmChartInstall.upgrade = false

	} else if err != nil {
		return err
	}

	return installChartAsync(helmChartInstall)

}
