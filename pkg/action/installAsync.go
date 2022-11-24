package action

import (
	"github.com/redhat-cop/oc-helm/pkg/options"
)

type InstallAsyncAction struct {
	baseAction
}

func NewInstallAsyncAction(commandLineOptions *options.CommandLineOption) *InstallAsyncAction {
	return &InstallAsyncAction{
		baseAction: baseAction{
			commandLineOptions: commandLineOptions,
		},
	}
}

func (i *InstallAsyncAction) Run(releaseName string, chartReference string) error {

	helmChartInstall := &HelmChartInstall{
		releaseName:        releaseName,
		chartReference:     chartReference,
		commandLineOptions: i.commandLineOptions,
		helmChartClient:    i.helmChartClient,
	}

	return installChartAsync(helmChartInstall)

}
