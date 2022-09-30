package action

import (
	"fmt"
	"text/tabwriter"

	"github.com/redhat-cop/oc-helm/pkg/client"
	"github.com/redhat-cop/oc-helm/pkg/options"
	"k8s.io/helm/pkg/strvals"
)

type VerifyAction struct {
	baseAction
}
type HelmVerify struct {
	commandLineOptions *options.CommandLineOption
	helmChartClient    *client.HelmChartClient
}

func NewVerifyAction(commandLineOptions *options.CommandLineOption) *VerifyAction {
	return &VerifyAction{
		baseAction: baseAction{
			commandLineOptions: commandLineOptions,
		},
	}
}

func (i *VerifyAction) Run() error {
	values, err := getValuesFromVerifyOptions(i.commandLineOptions.VerifierOptions)
	if err != nil {
		return err
	}
	result, err := i.helmChartClient.VerifyChart(i.commandLineOptions.ChartUrl, values)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(i.commandLineOptions.Streams.Out, 0, 8, 1, '\t', tabwriter.AlignRight)

	fmt.Fprintln(w, "NUMBER OF CHECKS PASSED:", result.VerifierApiResult.Passed)
	fmt.Fprintln(w, "NUMBER OF CHECKS FAILED:", result.VerifierApiResult.Failed)
	for _, message := range result.VerifierApiResult.Messages {
		fmt.Fprintf(w, "%s\t%s", "*", message)
		fmt.Fprint(w, "\n")
	}
	w.Flush()
	return nil

}

func getValuesFromVerifyOptions(verifierOptions []string) (map[string]interface{}, error) {
	values := make(map[string]interface{}, 0)
	for _, value := range verifierOptions {
		if err := strvals.ParseInto(value, values); err != nil {
			return nil, err
		}
	}
	return values, nil
}
