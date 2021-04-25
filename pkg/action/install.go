package action

import (
	"fmt"
	"strings"
	"time"

	"github.com/redhat-cop/oc-helm/pkg/client"
	"helm.sh/helm/v3/pkg/repo"

	"github.com/redhat-cop/oc-helm/pkg/options"
)

type InstallAction struct {
	commandLineOptions *options.CommandLineOption
	helmChartClient    *client.HelmChartClient
}

func NewInstallAction(commandLineOptions *options.CommandLineOption) *InstallAction {
	return &InstallAction{
		commandLineOptions: commandLineOptions,
	}
}

func (i *InstallAction) BuildHelmChartClient() error {

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

func (i *InstallAction) Run(releaseName string, chartReference string) error {

	index, err := i.helmChartClient.GetIndex()

	if err != nil {
		return err
	}

	chartReferenceParts := strings.Split(chartReference, "/")

	repository := chartReferenceParts[0]
	chartName := chartReferenceParts[1]

	consoleChartName := fmt.Sprintf("%s--%s", chartName, repository)

	charts := index.Entries[consoleChartName]

	if charts == nil {
		return fmt.Errorf("Error: Chart '%s' does not exist", chartReference)
	}

	chartVersion := i.getChartVersion(charts)

	if chartVersion == nil {

		if i.commandLineOptions.Version != "" {
			return fmt.Errorf("Chart '%s' with version '%s' not found", chartReference, i.commandLineOptions.Version)
		} else {
			return fmt.Errorf("Chart '%s' not found", chartReference)
		}

	}

	if len(chartVersion.URLs) == 0 {
		return fmt.Errorf("Unable to locate Chart URL")
	}

	chartURL := chartVersion.URLs[0]
	chart, err := i.helmChartClient.GetChart(chartURL)

	if err != nil {
		return err
	}

	// TODO: Manage Values

	release, err := i.helmChartClient.CreateRelease(releaseName, chartURL, chart.Values)

	if err != nil {
		return err
	}

	fmt.Fprintf(i.commandLineOptions.Streams.Out, "NAME: %s\n", release.Name)
	fmt.Fprintf(i.commandLineOptions.Streams.Out, "NAMESPACE: %s\n", release.Namespace)
	if !release.Info.LastDeployed.IsZero() {
		fmt.Fprintf(i.commandLineOptions.Streams.Out, "LAST DEPLOYED: %s\n", release.Info.LastDeployed.Format(time.ANSIC))
	}
	fmt.Fprintf(i.commandLineOptions.Streams.Out, "STATUS: %s\n", release.Info.Status.String())
	fmt.Fprintf(i.commandLineOptions.Streams.Out, "REVISION: %d\n", release.Version)

	return nil

}

func (i *InstallAction) getChartVersion(charts repo.ChartVersions) *repo.ChartVersion {

	if len(charts) == 0 {
		return nil
	}

	if i.commandLineOptions.Version == "" {
		return charts[0]
	}

	for _, chart := range charts {
		if chart.Version == i.commandLineOptions.Version {
			return chart
		}
	}

	return nil

}
