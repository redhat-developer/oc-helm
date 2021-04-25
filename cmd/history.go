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

	"github.com/redhat-cop/oc-helm/pkg/action"
	"github.com/redhat-cop/oc-helm/pkg/options"
	"github.com/spf13/cobra"
)

func newHistoryCmd(commandLineOptions *options.CommandLineOption) *cobra.Command {

	action := action.NewHistoryAction(commandLineOptions)

	// listCmd represents the list command
	historyCmd := &cobra.Command{
		Use:     "history RELEASE_NAME",
		Short:   "Fetch release history",
		Aliases: []string{"hist"},
		Long:    "Fetch release history",
		PreRunE: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("Error: \"RELEASE_NAME\" argument is required")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			err := action.BuildHelmChartClient()

			if err != nil {
				return err
			}

			return action.Run(args[0])

		},
	}

	return historyCmd
}
