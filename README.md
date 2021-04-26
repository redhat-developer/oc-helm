OpenShift Helm CLI tool (oc-helm)

This is a command line tool to execute helm commands via the OpenShift API/CLI.

## Build tool

```shell
go build
sudo cp oc-helm /usr/local/bin/.
```

## Usage

```shell
Usage:
  oc-helm [command]

Available Commands:
  help        Help about any command
  history     Fetch release history
  index       Index of Available Charts
  install     Install chart
  list        List installed charts
  rollback    Roll back a release to a previous revision
  uninstall   Uninstall a Release
  upgrade     Upgrade a release

Flags:
      --console-hostname string   OpenShift Console Hostname
      --context string            Kubernetes Context
  -h, --help                      help for oc-helm
  -n, --namespace string          Kubernetes namespace
  -t, --token string              OpenShift OAuth token

Use "oc-helm [command] --help" for more information about a command.
```

