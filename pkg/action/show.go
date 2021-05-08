package action

import (
	"fmt"
	"strings"

	"github.com/redhat-cop/oc-helm/pkg/options"
	"github.com/redhat-cop/oc-helm/pkg/utils"
	"helm.sh/helm/v3/pkg/chart"
	"sigs.k8s.io/yaml"
)

type ShowAction struct {
	baseAction
	ShowFormat ShowFormat
}

type ShowFormat string

const (
	ShowAll        ShowFormat = "all"
	ShowChart      ShowFormat = "chart"
	ShowValues     ShowFormat = "values"
	ShowReadme     ShowFormat = "readme"
	ValuesfileName string     = "values.yaml"
)

var readmeFileNames = []string{"readme.md", "readme.txt", "readme"}

func NewShowAction(commandLineOptions *options.CommandLineOption) *ShowAction {
	return &ShowAction{
		baseAction: baseAction{
			commandLineOptions: commandLineOptions,
		},
	}
}

func (s *ShowAction) Run(chartReference string) error {

	index, err := s.helmChartClient.GetIndex()

	if err != nil {
		return err
	}

	index.SortEntries()

	chartReferenceParts := strings.Split(chartReference, "/")

	repository := chartReferenceParts[0]
	chartName := chartReferenceParts[1]

	consoleChartName := utils.CreateRepositoryIndexKey(repository, chartName)

	charts := index.Entries[consoleChartName]

	if charts == nil {
		return fmt.Errorf("Error: Chart '%s' does not exist", chartReference)
	}

	chartVersion := getChartVersion(s.commandLineOptions, charts)

	if chartVersion == nil {

		if s.commandLineOptions.Version != "" {
			return fmt.Errorf("Chart '%s' with version '%s' not found", chartReference, s.commandLineOptions.Version)
		} else {
			return fmt.Errorf("Chart '%s' not found", chartReference)
		}

	}

	if len(chartVersion.URLs) == 0 {
		return fmt.Errorf("Unable to locate Chart URL")
	}

	chart, err := s.helmChartClient.GetChart(chartVersion.URLs[0])

	if err != nil {
		return fmt.Errorf("Could not obtain chart '%s' with url: '%s'", chartReference, chartVersion.URLs[0])
	}

	var out strings.Builder

	cf, err := yaml.Marshal(chart.Metadata)
	if err != nil {
		return err
	}

	if s.ShowFormat == ShowChart || s.ShowFormat == ShowAll {
		fmt.Fprintf(&out, "%s\n", cf)
	}

	if (s.ShowFormat == ShowValues || s.ShowFormat == ShowAll) && chart.Values != nil {
		if s.ShowFormat == ShowAll {
			fmt.Fprintln(&out, "---")
		}

		values, err := yaml.Marshal(chart.Values)

		if err == nil {
			fmt.Fprintln(&out, string(values))
		}
	}

	if s.ShowFormat == ShowReadme || s.ShowFormat == ShowAll {
		if s.ShowFormat == ShowAll {
			fmt.Fprintln(&out, "---")
		}
		readme := findReadme(chart.Files)
		if readme != nil {
			fmt.Fprintf(&out, "%s\n", readme.Data)
		}
	}

	fmt.Fprintln(s.commandLineOptions.Streams.Out, out.String())

	return nil
}

func findReadme(files []*chart.File) (file *chart.File) {
	for _, file := range files {
		for _, n := range readmeFileNames {
			if strings.EqualFold(file.Name, n) {
				return file
			}
		}
	}
	return nil
}
