package action

import (
	"github.com/redhat-cop/oc-helm/pkg/options"
)

type InstallAction struct {
	baseAction
}

func NewInstallAction(commandLineOptions *options.CommandLineOption) *InstallAction {
	return &InstallAction{
		baseAction: baseAction{
			commandLineOptions: commandLineOptions,
		},
	}
}

func (i *InstallAction) Run(releaseName string, chartReference string) error {

	helmChartInstall := &HelmChartInstall{
		releaseName:        releaseName,
		chartReference:     chartReference,
		commandLineOptions: i.commandLineOptions,
		helmChartClient:    i.helmChartClient,
	}

	return installChart(helmChartInstall)

}
