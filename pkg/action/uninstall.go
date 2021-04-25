package action

import (
	"fmt"

	"github.com/redhat-cop/oc-helm/pkg/client"

	"github.com/redhat-cop/oc-helm/pkg/options"
)

type UninstallAction struct {
	commandLineOptions *options.CommandLineOption
	helmChartClient    *client.HelmChartClient
}

func NewUninstallAction(commandLineOptions *options.CommandLineOption) *UninstallAction {
	return &UninstallAction{
		commandLineOptions: commandLineOptions,
	}
}

func (h *UninstallAction) BuildHelmChartClient() error {

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

func (u *UninstallAction) Run(releaseName string) error {

	release, err := u.helmChartClient.Uninstall(releaseName)

	if err != nil {
		return err
	}

	if release != nil && release.Info != "" {
		fmt.Fprintln(u.commandLineOptions.Streams.Out, release.Info)
	} else {
		fmt.Fprintf(u.commandLineOptions.Streams.Out, "release \"%s\" uninstalled\n", releaseName)
	}

	return nil

}
