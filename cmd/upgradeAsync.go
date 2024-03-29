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
	"fmt"
	"strings"

	"github.com/redhat-cop/oc-helm/pkg/action"
	"github.com/redhat-cop/oc-helm/pkg/options"
	"github.com/spf13/cobra"
)

func newUpgradeAsyncCmd(commandLineOptions *options.CommandLineOption) *cobra.Command {

	action := action.NewUpgradeAsyncAction(commandLineOptions)

	// upgradeCmd represents the upgrade command
	upgradeCmd := &cobra.Command{
		Use:   "upgrade-async [NAME] [REPOSITORY/NAME]",
		Short: "Upgrade a release",
		Long:  "Upgrade a release",
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

			err := action.Setup()

			if err != nil {
				return err
			}

			return action.Run(args[0], args[1])
		},
	}

	upgradeCmd.PersistentFlags().StringVar(&commandLineOptions.Version, "version", "", "specify the exact chart version to use. If this is not specified, the latest version is used")
	upgradeCmd.PersistentFlags().BoolVarP(&commandLineOptions.Install, "install", "i", false, "if a release by this name doesn't already exist, run an install")
	setValuesOptions(upgradeCmd, commandLineOptions)

	return upgradeCmd
}
