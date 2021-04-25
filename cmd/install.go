/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"fmt"
	"strings"

	"github.com/redhat-cop/oc-helm/pkg/action"
	"github.com/redhat-cop/oc-helm/pkg/options"
	"github.com/spf13/cobra"
)

func newInstallCmd(commandLineOptions *options.CommandLineOption) *cobra.Command {

	action := action.NewInstallAction(commandLineOptions)

	// installCmd represents the install command
	installCmd := &cobra.Command{
		Use:   "install [NAME] [REPOSITORY/NAME]",
		Short: "Install chart",
		Long:  `Install chart`,
		PreRunE: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("Error: Must provide a release name and reference to chart within a repository")
			}

			// Check format of chart reference
			chartReferenceParts := strings.Split(args[1], "/")

			if len(chartReferenceParts) != 2 {
				return fmt.Errorf("Error: Chart reference must take the form '<REPOSITORY/NAME>'")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			err := action.BuildHelmChartClient()

			if err != nil {
				return err
			}

			return action.Run(args[0], args[1])

		},
	}

	installCmd.PersistentFlags().StringVar(&commandLineOptions.Version, "version", "", "specify the exact chart version to use. If this is not specified, the latest version is used")

	return installCmd
}
