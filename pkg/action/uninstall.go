package action

import (
	"fmt"

	"github.com/redhat-cop/oc-helm/pkg/options"
)

type UninstallAction struct {
	baseAction
}

func NewUninstallAction(commandLineOptions *options.CommandLineOption) *UninstallAction {
	return &UninstallAction{
		baseAction: baseAction{
			commandLineOptions: commandLineOptions,
		},
	}
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
