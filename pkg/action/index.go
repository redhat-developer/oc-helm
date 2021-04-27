package action

import (
	"fmt"
	"sort"
	"text/tabwriter"

	"github.com/redhat-cop/oc-helm/pkg/types"
	"github.com/redhat-cop/oc-helm/pkg/utils"

	"helm.sh/helm/v3/pkg/repo"

	"github.com/redhat-cop/oc-helm/pkg/options"
)

type IndexAction struct {
	baseAction
}

func NewIndexAction(commandLineOptions *options.CommandLineOption) *IndexAction {
	return &IndexAction{
		baseAction: baseAction{
			commandLineOptions: commandLineOptions,
		},
	}
}

func (i *IndexAction) Run() error {

	index, err := i.helmChartClient.GetIndex()

	if err != nil {
		return err
	}

	if len(index.Entries) > 0 {

		index.SortEntries()

		charVersionRepositories := sortIndexEntries(index.Entries)

		w := tabwriter.NewWriter(i.commandLineOptions.Streams.Out, 0, 8, 1, '\t', tabwriter.AlignRight)

		fmt.Fprintln(w, "REPOSITORY\tNAME\tLATEST VERSION")

		for _, chartVersionRepository := range charVersionRepositories {

			fmt.Fprintf(w, "%s\t%s\t", chartVersionRepository.Repository, chartVersionRepository.Chart)

			if len(chartVersionRepository.ChartVersions) > 0 {
				fmt.Fprint(w, chartVersionRepository.ChartVersions[0].Version)
			}

			fmt.Fprint(w, "\n")

		}

		w.Flush()

	} else {
		fmt.Fprintln(i.commandLineOptions.Streams.Out, "No Charts Found")
	}

	return nil

}

func sortIndexEntries(entries map[string]repo.ChartVersions) []types.ChartVersionRepository {

	uniqueRepositories := map[string][]string{}
	repositories := []string{}
	chartVersions := []types.ChartVersionRepository{}

	for key, _ := range entries {
		repository, chart, err := utils.SplitRepositoryIndexKey(key)
		if err == nil {
			charts := uniqueRepositories[repository]

			if charts == nil {
				repositories = append(repositories, repository)
				charts = []string{chart}
			} else {
				charts = append(charts, chart)
			}

			uniqueRepositories[repository] = charts

		}
	}

	sort.Strings(repositories)

	for _, repositoryKey := range repositories {

		repositoryCharts := uniqueRepositories[repositoryKey]

		sort.Strings(repositoryCharts)

		for _, repositoryChart := range repositoryCharts {

			if repositoryChartVersions, ok := entries[utils.CreateRepositoryIndexKey(repositoryKey, repositoryChart)]; ok {
				chartVersions = append(chartVersions, types.ChartVersionRepository{
					Chart:         repositoryChart,
					Repository:    repositoryKey,
					ChartVersions: repositoryChartVersions,
				})
			}

		}
	}

	return chartVersions
}
