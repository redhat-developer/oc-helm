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

	"github.com/redhat-cop/oc-helm/pkg/action"
	"github.com/redhat-cop/oc-helm/pkg/options"
	"github.com/spf13/cobra"
)

func newUninstallAsyncCmd(commandLineOptions *options.CommandLineOption) *cobra.Command {

	action := action.NewUninstallActionAsync(commandLineOptions)

	// uninstallCmd represents the uninstall command
	uninstallCmd := &cobra.Command{
		Use:     "uninstall-async RELEASE_NAME REVISION",
		Short:   "Uninstall a Release asynchronously",
		Aliases: []string{"del", "delete", "un"},
		Long:    "Uninstall a Release asynchronously",
		PreRunE: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("Error: \"RELEASE_NAME\" argument is required and \"REVISION\" argument is required")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			err := action.Setup()

			if err != nil {
				return err
			}

			err = action.Run(args[0], args[1])
			if err != nil {
				return err
			}
			return nil
		},
	}

	uninstallCmd.PersistentFlags().StringVar(&commandLineOptions.Revision, "revision", "-r", "specify the exact release revision to uninstall. If this is not specified, the latest revision will be deleted")
	return uninstallCmd
}
