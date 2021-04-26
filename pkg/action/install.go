package action

import (
	"github.com/redhat-cop/oc-helm/pkg/client"

	"github.com/redhat-cop/oc-helm/pkg/options"
)

type InstallAction struct {
	commandLineOptions *options.CommandLineOption
	helmChartClient    *client.HelmChartClient
}

func NewInstallAction(commandLineOptions *options.CommandLineOption) *InstallAction {
	return &InstallAction{
		commandLineOptions: commandLineOptions,
	}
}

func (i *InstallAction) BuildHelmChartClient() error {

	if err := i.commandLineOptions.Process(); err != nil {
		return err
	}

	helmChartClient, err := client.NewHelmChartClient(i.commandLineOptions)

	if err != nil {
		return err
	}

	i.helmChartClient = helmChartClient

	return nil

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
