/*
Copyright © 2022 redhat-developer

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/redhat-cop/oc-helm/pkg/action"
	"github.com/redhat-cop/oc-helm/pkg/options"
	"github.com/spf13/cobra"
	"helm.sh/helm/v3/cmd/helm/require"
)

func newShowCmd(commandLineOptions *options.CommandLineOption) *cobra.Command {

	showAction := action.NewShowAction(commandLineOptions)

	// showCmd represents the show command
	showCmd := &cobra.Command{
		Use:   "show",
		Short: "show information of a chart",
		Long:  "show information of a chart",
		Args:  require.NoArgs,
	}

	allCmd := &cobra.Command{
		Use:   "all [CHART]",
		Short: "show all information of the chart",
		Long:  "show all information of the chart",
		Args:  require.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			showAction.ShowFormat = action.ShowAll

			err := showAction.Setup()

			if err != nil {
				return err
			}

			return showAction.Run(args[0])

		},
	}

	valuesCmd := &cobra.Command{
		Use:   "values [CHART]",
		Short: "show the chart's values",
		Long:  "show the chart's values",
		Args:  require.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			showAction.ShowFormat = action.ShowValues

			err := showAction.Setup()

			if err != nil {
				return err
			}

			return showAction.Run(args[0])

		},
	}

	chartCmd := &cobra.Command{
		Use:   "chart [CHART]",
		Short: "show the chart's definition",
		Long:  "show the chart's definition",
		Args:  require.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			showAction.ShowFormat = action.ShowChart

			err := showAction.Setup()

			if err != nil {
				return err
			}

			return showAction.Run(args[0])

		},
	}

	readmeCmd := &cobra.Command{
		Use:   "readme [CHART]",
		Short: "show the chart's README",
		Long:  "show the chart's README",
		Args:  require.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			showAction.ShowFormat = action.ShowReadme

			err := showAction.Setup()

			if err != nil {
				return err
			}

			return showAction.Run(args[0])

		},
	}

	cmds := []*cobra.Command{allCmd, readmeCmd, valuesCmd, chartCmd}
	for _, subCmd := range cmds {
		subCmd.PersistentFlags().StringVar(&commandLineOptions.Version, "version", "", "specify the exact chart version to use. If this is not specified, the latest version is used")
		showCmd.AddCommand(subCmd)
	}

	return showCmd
}
