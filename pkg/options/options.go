package options

import (
	"fmt"
	"os"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type CommandLineOption struct {
	ConfigFlags     *genericclioptions.ConfigFlags
	Streams         *genericclioptions.IOStreams
	ConsoleHostname string
	Token           string
	Insecure        bool
	Namespace       string
	APIServer       string
	ChartName       string
	ChartRepository string
	ChartUrl        string
	Context         string
	Version         string
	ValueFiles      []string
	StringValues    []string
	Values          []string
	VerifierOptions []string
	FileValues      []string
	Install         bool
}

func NewCommandLineOption() *CommandLineOption {

	return &CommandLineOption{
		ConfigFlags: genericclioptions.NewConfigFlags(true),
		Streams: &genericclioptions.IOStreams{
			In:     os.Stdin,
			Out:    os.Stdout,
			ErrOut: os.Stderr,
		}}

}

func (c *CommandLineOption) Process() error {
	if c.Token != "" && c.ConsoleHostname != "" {
		return nil
	}

	rawConfig, err := c.ConfigFlags.ToRawKubeConfigLoader().RawConfig()

	if err != nil {
		return nil
	}

	if c.Context == "" {
		c.Context = rawConfig.CurrentContext
	}

	currentContext, exists := rawConfig.Contexts[c.Context]

	if !exists {
		return fmt.Errorf("Error: No Current Context Exists")
	}

	if c.Token == "" {
		c.Token = rawConfig.AuthInfos[currentContext.AuthInfo].Token
	}

	if c.Namespace == "" {
		c.Namespace = currentContext.Namespace
	}

	if c.APIServer == "" {
		c.APIServer = rawConfig.Clusters[currentContext.Cluster].Server
	}

	return nil

}
