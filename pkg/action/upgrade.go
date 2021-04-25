package action

import (
	"github.com/redhat-cop/oc-helm/pkg/client"

	"github.com/redhat-cop/oc-helm/pkg/options"
)

type UpgradeAction struct {
	commandLineOptions *options.CommandLineOption
	helmChartClient    *client.HelmChartClient
}

func NewUpgradeAction(commandLineOptions *options.CommandLineOption) *UpgradeAction {
	return &UpgradeAction{
		commandLineOptions: commandLineOptions,
	}
}

func (h *UpgradeAction) BuildHelmChartClient() error {

	if err := h.commandLineOptions.Process(); err != nil {
		return err
	}

	helmChartClient, err := client.NewHelmChartClient(h.commandLineOptions)

	if err != nil {
		return err
	}

	h.helmChartClient = helmChartClient

	return nil

}

func (i *UpgradeAction) Run(releaseName string, chartReference string) error {

	return nil

}
