/*
Copyright © 2022  redhat-developer

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
	"text/template"

	"github.com/redhat-cop/oc-helm/pkg/options"
	"github.com/redhat-cop/oc-helm/pkg/version"
	"github.com/spf13/cobra"
)

const (
	versionTemplate = `{{.Version}}
`
	fullInfoTemplate = `oc-helm:    {{.Version}}
platform:   {{.Platform}}
git commit: {{.GitCommit}}
build date: {{.BuildDate}}
go version: {{.GoVersion}}
compiler:   {{.Compiler}}
`
	flagFull = "full"
)

func newVersionCmd(commandLineOptions *options.CommandLineOption) *cobra.Command {

	// indexCmd represents the index command
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Long:  "Print the version information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {

			var tpl string

			if cmd.Flag(flagFull).Changed {
				tpl = fullInfoTemplate
			} else {
				tpl = versionTemplate
			}

			var t = template.Must(template.New("info").Parse(tpl))

			if err := t.Execute(commandLineOptions.Streams.Out, version.GetBuildInfo()); err != nil {
				return fmt.Errorf("Could not print version info")
			}

			return nil

		},
	}

	versionCmd.Flags().BoolP(flagFull, "f", false, "print extended version information")

	return versionCmd

}
