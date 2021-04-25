package action

import (
	"fmt"

	"github.com/redhat-cop/oc-helm/pkg/client"

	"github.com/redhat-cop/oc-helm/pkg/options"
)

type RollbackAction struct {
	commandLineOptions *options.CommandLineOption
	helmChartClient    *client.HelmChartClient
	revision           int
}

func NewRollbackAction(commandLineOptions *options.CommandLineOption) *RollbackAction {
	return &RollbackAction{
		commandLineOptions: commandLineOptions,
	}
}

func (r *RollbackAction) BuildHelmChartClient() error {

	if err := r.commandLineOptions.Process(); err != nil {
		return err
	}

	helmChartClient, err := client.NewHelmChartClient(r.commandLineOptions)

	if err != nil {
		return err
	}

	r.helmChartClient = helmChartClient

	return nil

}

func (r *RollbackAction) SetRevision(revision int) {
	r.revision = revision
}

func (r *RollbackAction) Run(releaseName string) error {

	if r.revision <= 0 {

		revision, err := r.helmChartClient.GetRelease(releaseName)

		if err != nil {
			return err
		}

		r.revision = revision.Version - 1

	}

	if r.revision < 1 {
		return fmt.Errorf("Error: release: not found")
	}

	_, err := r.helmChartClient.Rollback(releaseName, r.revision)

	if err != nil {
		return err
	}

	fmt.Fprintf(r.commandLineOptions.Streams.Out, "Rollback was a success! Happy Helming!\n")

	return nil

}
