package action

import (
	"fmt"

	"github.com/redhat-cop/oc-helm/pkg/options"
)

type UninstallActionAsync struct {
	baseAction
}

func NewUninstallActionAsync(commandLineOptions *options.CommandLineOption) *UninstallActionAsync {
	return &UninstallActionAsync{
		baseAction: baseAction{
			commandLineOptions: commandLineOptions,
		},
	}
}

func (u *UninstallActionAsync) Run(releaseName, version string) error {

	statusCode, err := u.helmChartClient.UninstallAsync(releaseName, version)

	if statusCode == -1 {
		return err
	}

	fmt.Fprintf(u.commandLineOptions.Streams.Out, "STATUS CODE: %v\n", statusCode)
	fmt.Fprintf(u.commandLineOptions.Streams.Out, "ERROR MESSAGE: %s\n", err)
	return nil

}
