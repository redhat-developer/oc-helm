package action

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/redhat-cop/oc-helm/pkg/client"

	"github.com/redhat-cop/oc-helm/pkg/options"
)

type IndexAction struct {
	commandLineOptions *options.CommandLineOption
	helmChartClient    *client.HelmChartClient
}

func NewIndexAction(commandLineOptions *options.CommandLineOption) Action {
	return &IndexAction{
		commandLineOptions: commandLineOptions,
	}
}

func (i *IndexAction) BuildHelmChartClient() error {

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

func (i *IndexAction) Run() error {

	index, err := i.helmChartClient.GetIndex()

	if err != nil {
		return err
	}

	if len(index.Entries) > 0 {

		index.SortEntries()

		w := tabwriter.NewWriter(i.commandLineOptions.Streams.Out, 0, 8, 1, '\t', tabwriter.AlignRight)

		fmt.Fprintln(w, "REPOSITORY\tNAME\tLATEST VERSION")

		for chartName, charts := range index.Entries {

			chartNameItems := strings.Split(chartName, "--")

			fmt.Fprintf(w, "%s\t%s\t", chartNameItems[1], chartNameItems[0])

			if len(charts) > 0 {
				fmt.Fprint(w, charts[0].Version)
			}

			fmt.Fprint(w, "\n")

		}

		w.Flush()

	} else {
		fmt.Fprintln(i.commandLineOptions.Streams.Out, "No Charts Found")
	}

	return nil

}
