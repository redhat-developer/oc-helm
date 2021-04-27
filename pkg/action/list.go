package action

import (
	"fmt"
	"text/tabwriter"
	"time"

	"github.com/redhat-cop/oc-helm/pkg/options"
)

type ListAction struct {
	baseAction
}

func NewListAction(commandLineOptions *options.CommandLineOption) *ListAction {
	return &ListAction{
		baseAction: baseAction{
			commandLineOptions: commandLineOptions,
		},
	}
}

func (l *ListAction) Run() error {

	releases, err := l.helmChartClient.ListReleases()

	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(l.commandLineOptions.Streams.Out, 0, 8, 1, '\t', tabwriter.AlignRight)

	fmt.Fprintln(w, "NAME\tNAMESPACE\tREVISION\tUPDATED\tSTATUS\tCHART\tAPP VERSION")

	for _, release := range *releases {

		fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\t%s\t%s", release.Name, release.Namespace, release.Version, release.Info.LastDeployed.Format(time.ANSIC), release.Info.Status, fmt.Sprintf("%s-%s", release.Chart.Metadata.Name, release.Chart.Metadata.Version), release.Chart.AppVersion())

		fmt.Fprint(w, "\n")

	}

	w.Flush()

	return nil

}
