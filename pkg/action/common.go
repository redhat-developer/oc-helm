package action

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/redhat-cop/oc-helm/pkg/client"
	"github.com/redhat-cop/oc-helm/pkg/options"
	"github.com/redhat-cop/oc-helm/pkg/utils"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
	kv1 "k8s.io/api/core/v1"
)

type HelmChartInstall struct {
	releaseName        string
	chartReference     string
	commandLineOptions *options.CommandLineOption
	helmChartClient    *client.HelmChartClient
	upgrade            bool
}

func getChartVersion(commandLineOptions *options.CommandLineOption, charts repo.ChartVersions) *repo.ChartVersion {

	if len(charts) == 0 {
		return nil
	}

	if commandLineOptions.Version == "" {
		return charts[0]
	}

	for _, chart := range charts {
		if chart.Version == commandLineOptions.Version {
			return chart
		}
	}

	return nil

}

func installChart(helmChartInstall *HelmChartInstall) error {
	index, err := helmChartInstall.helmChartClient.GetIndex()

	if err != nil {
		return err
	}

	index.SortEntries()

	chartReferenceParts := strings.Split(helmChartInstall.chartReference, "/")

	repository := chartReferenceParts[0]
	chartName := chartReferenceParts[1]

	consoleChartName := utils.CreateRepositoryIndexKey(repository, chartName)

	charts := index.Entries[consoleChartName]

	if charts == nil {
		return fmt.Errorf("Error: Chart '%s' does not exist", helmChartInstall.chartReference)
	}

	chartVersion := getChartVersion(helmChartInstall.commandLineOptions, charts)

	if chartVersion == nil {

		if helmChartInstall.commandLineOptions.Version != "" {
			return fmt.Errorf("Chart '%s' with version '%s' not found", helmChartInstall.chartReference, helmChartInstall.commandLineOptions.Version)
		} else {
			return fmt.Errorf("Chart '%s' not found", helmChartInstall.chartReference)
		}

	}

	if len(chartVersion.URLs) == 0 {
		return fmt.Errorf("Unable to locate Chart URL")
	}

	chartURL := chartVersion.URLs[0]

	values, err := utils.MergeValues(helmChartInstall.commandLineOptions)

	if err != nil {
		return err
	}

	release, err := helmChartInstall.helmChartClient.CreateRelease(helmChartInstall.releaseName, chartURL, values, helmChartInstall.upgrade)

	if err != nil {
		return err
	}

	printReleaseDeploymentStatus(helmChartInstall.commandLineOptions.Streams.Out, release)

	return nil
}

func printReleaseDeploymentStatus(w io.Writer, release *release.Release) {
	fmt.Fprintf(w, "NAME: %s\n", release.Name)
	fmt.Fprintf(w, "NAMESPACE: %s\n", release.Namespace)
	if !release.Info.LastDeployed.IsZero() {
		fmt.Fprintf(w, "LAST DEPLOYED: %s\n", release.Info.LastDeployed.Format(time.ANSIC))
	}
	fmt.Fprintf(w, "STATUS: %s\n", release.Info.Status.String())
	fmt.Fprintf(w, "REVISION: %d\n", release.Version)

}

func installChartAsync(helmChartInstall *HelmChartInstall) error {
	index, err := helmChartInstall.helmChartClient.GetIndex()

	if err != nil {
		return err
	}

	index.SortEntries()

	chartReferenceParts := strings.Split(helmChartInstall.chartReference, "/")

	repository := chartReferenceParts[0]
	chartName := chartReferenceParts[1]

	consoleChartName := utils.CreateRepositoryIndexKey(repository, chartName)

	charts := index.Entries[consoleChartName]

	if charts == nil {
		return fmt.Errorf("Error: Chart '%s' does not exist", helmChartInstall.chartReference)
	}

	chartVersion := getChartVersion(helmChartInstall.commandLineOptions, charts)

	if chartVersion == nil {

		if helmChartInstall.commandLineOptions.Version != "" {
			return fmt.Errorf("Chart '%s' with version '%s' not found", helmChartInstall.chartReference, helmChartInstall.commandLineOptions.Version)
		} else {
			return fmt.Errorf("Chart '%s' not found", helmChartInstall.chartReference)
		}

	}

	if len(chartVersion.URLs) == 0 {
		return fmt.Errorf("Unable to locate Chart URL")
	}

	chartURL := chartVersion.URLs[0]

	values, err := utils.MergeValues(helmChartInstall.commandLineOptions)

	if err != nil {
		return err
	}

	release, err := helmChartInstall.helmChartClient.CreateReleaseAsync(helmChartInstall.releaseName, chartURL, values, helmChartInstall.upgrade)

	if err != nil {
		return err
	}

	printReleaseDeploymentStatusAsync(helmChartInstall.commandLineOptions.Streams.Out, release)

	return nil
}

func printReleaseDeploymentStatusAsync(w io.Writer, secret *kv1.Secret) {
	fmt.Fprintf(w, "NAME: %s\n", secret.ObjectMeta.Name)
	fmt.Fprintf(w, "UID: %s\n", secret.ObjectMeta.UID)
}
