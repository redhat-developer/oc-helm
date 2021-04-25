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
	"os"

	"github.com/spf13/cobra"

	"github.com/redhat-cop/oc-helm/pkg/options"
)

func newRootCmd(commandLineOptions *options.CommandLineOption) (*cobra.Command, error) {

	rootCmd := &cobra.Command{
		Use:           "oc-helm",
		Short:         "helm is oc plugin that integrates with Helm and OpenShift",
		Long:          "OpenShift Command Line tool to interact with Helm capabilities.",
		SilenceUsage:  true,
		SilenceErrors: true,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
	}

	rootCmd.PersistentFlags().StringVar(&commandLineOptions.ConsoleHostname, "console-hostname", "", "OpenShift Console Hostname")
	rootCmd.PersistentFlags().StringVar(&commandLineOptions.Context, "context", "", "Kubernetes Context")
	rootCmd.PersistentFlags().StringVarP(&commandLineOptions.Namespace, "namespace", "n", "", "Kubernetes namespace")
	rootCmd.PersistentFlags().StringVarP(&commandLineOptions.Token, "token", "t", "", "OpenShift OAuth token")

	rootCmd.AddCommand(
		newHistoryCmd(commandLineOptions),
		newIndexCmd(commandLineOptions),
		newInstallCmd(commandLineOptions),
		newListCmd(commandLineOptions),
		newRollbackCmd(commandLineOptions),
		newUninstallCmd(commandLineOptions),
	)

	return rootCmd, nil

}

func Execute() {

	commandLineOption := options.NewCommandLineOption()

	rootCmd, _ := newRootCmd(commandLineOption)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(commandLineOption.Streams.ErrOut, err)
		os.Exit(1)
	}
}
