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
	"os"

	"github.com/redhat-cop/oc-helm/pkg/constants"
	"github.com/redhat-cop/oc-helm/pkg/options"
	"github.com/spf13/cobra"
)

func newRootCmd(commandLineOptions *options.CommandLineOption) (*cobra.Command, error) {

	rootCmd := &cobra.Command{
		Use:           "oc-helm",
		Short:         "helm is oc plugin that integrates with Helm and OpenShift",
		Long:          "OpenShift Command Line tool to interact with Helm capabilities.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.PersistentFlags().StringVar(&commandLineOptions.APIServer, "apiserver", os.Getenv(constants.OPENSHIFT_APISERVER), "OpenShift API Server Hostname")
	rootCmd.PersistentFlags().StringVar(&commandLineOptions.ConsoleHostname, "console-hostname", os.Getenv(constants.OPENSHIFT_CONSOLE_HOSTNAME_ENV), "OpenShift Console Hostname")
	rootCmd.PersistentFlags().StringVar(&commandLineOptions.Context, "context", "", "Kubernetes Context")
	rootCmd.PersistentFlags().StringVarP(&commandLineOptions.Namespace, "namespace", "n", "", "Kubernetes namespace")
	rootCmd.PersistentFlags().StringVarP(&commandLineOptions.Token, "token", "t", os.Getenv(constants.OPENSHIFT_OAUTH_TOKEN_ENV), "OpenShift OAuth token")

	rootCmd.AddCommand(
		newHistoryCmd(commandLineOptions),
		newIndexCmd(commandLineOptions),
		newInstallCmd(commandLineOptions),
		newListCmd(commandLineOptions),
		newRollbackCmd(commandLineOptions),
		newUninstallCmd(commandLineOptions),
		newUpgradeCmd(commandLineOptions),
		newVersionCmd(commandLineOptions),
		newShowCmd(commandLineOptions),
		newVerifyCmd(commandLineOptions),
		newInstallAsyncCmd(commandLineOptions),
		newUpgradeAsyncCmd(commandLineOptions),
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

func setValuesOptions(cmd *cobra.Command, commandLineOption *options.CommandLineOption) {
	cmd.Flags().StringSliceVarP(&commandLineOption.ValueFiles, "values", "f", []string{}, "specify values in a YAML file or a URL (can specify multiple)")
	cmd.Flags().StringArrayVar(&commandLineOption.Values, "set", []string{}, "set values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")
	cmd.Flags().StringArrayVar(&commandLineOption.StringValues, "set-string", []string{}, "set STRING values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")
	cmd.Flags().StringArrayVar(&commandLineOption.FileValues, "set-file", []string{}, "set values from respective files specified via the command line (can specify multiple or separate values with commas: key1=path1,key2=path2)")

}
