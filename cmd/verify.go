package cmd

import (
	"github.com/redhat-cop/oc-helm/pkg/action"
	"github.com/redhat-cop/oc-helm/pkg/options"
	"github.com/spf13/cobra"
)

// NewVerifyCmd creates ...
func newVerifyCmd(commandLineOptions *options.CommandLineOption) *cobra.Command {

	action := action.NewVerifyAction(commandLineOptions)

	verifyCmd := &cobra.Command{
		Use:   "verify",
		Short: "Verifies a Helm chart by checking some of its characteristics",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := action.Setup()
			if err != nil {
				return err
			}
			return action.Run()
		},
	}
	// settings.AddFlags(verifyCmd.Flags())
	verifyCmd.Flags().StringVarP(&commandLineOptions.ChartUrl, "chart-url", "", "", "chart url of the chart to be verified")
	verifyCmd.Flags().StringSliceVarP(&commandLineOptions.VerifierOptions, "values", "V", nil, "set the profile with which the chart url is to be validated (example:provider=developer-console)")
	verifyCmd.MarkPersistentFlagRequired("chart-url")
	return verifyCmd
}
