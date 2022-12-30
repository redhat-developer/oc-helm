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
)

func newListCmd(commandLineOptions *options.CommandLineOption) *cobra.Command {

	action := action.NewListAction(commandLineOptions)

	// listCmd represents the list command
	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List installed charts",
		Long:    `List installed charts`,
		RunE: func(cmd *cobra.Command, args []string) error {

			err := action.Setup()

			if err != nil {
				return err
			}

			return action.Run()

		},
	}
	listCmd.Flags().BoolVarP(&commandLineOptions.LimitInfo, "limitInfo", "i", false, "specifies if call is made from topology or not")
	return listCmd
}
