package action

import (
	"fmt"
	"text/tabwriter"
	"time"

	"github.com/redhat-cop/oc-helm/pkg/options"
)

type HistoryAction struct {
	baseAction
}

func NewHistoryAction(commandLineOptions *options.CommandLineOption) *HistoryAction {
	return &HistoryAction{
		baseAction: baseAction{
			commandLineOptions: commandLineOptions,
		},
	}
}

func (h *HistoryAction) Run(releaseName string) error {

	releases, err := h.helmChartClient.History(releaseName)

	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(h.commandLineOptions.Streams.Out, 0, 8, 1, '\t', tabwriter.AlignRight)

	fmt.Fprintln(w, "REVISION\tUPDATED\tSTATUS\tCHART\tAPP VERSION\tDESCRIPTION")

	for _, release := range *releases {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s", release.Version, release.Info.LastDeployed.Format(time.ANSIC), release.Info.Status, fmt.Sprintf("%s-%s", release.Chart.Metadata.Name, release.Chart.Metadata.Version), release.Chart.AppVersion(), release.Info.Description)

		fmt.Fprint(w, "\n")

	}

	w.Flush()

	return nil

}
