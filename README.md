# oc helm

OpenShift CLI plugin to integrate with platform capabilities supporting [Helm](https://helm.sh/).

## Overview

OpenShift provides support for managing the lifecycle of Helm charts. This capability is limited primarily to the [Web Console](https://docs.openshift.com/container-platform/4.7/applications/application_life_cycle_management/odc-working-with-helm-charts-using-developer-perspective.html). This plugin enables the management of Helm charts similar to using the standalone Helm CLI while offloading much of the work to OpenShift.

## Capabilities

The following capabilities are provides by this plugin

* Discovering Helm Charts and Repositories registered on the platform
* Chart lifecycle
  * Show
  * Install
  * Upgrade
  * Rollback
  * History
  * List
  * Uninstall
  * Verify

## Prerequisites

The following prerequisites must be met prior to using the plugin:

1. OpenShift CLI
2. OpenShift environment
    1. You must be logged in using the OpenShift CLI or have a valid environment and OAuth token

## Installing

Perform the following steps to setup and configure the plugin on your machine:

1. Download the latest release for your operating system from the [Release Page](https://github.com/sabre1041/oc-helm/releases)

2. Extract the compressed archive and move the resulting binary to your path

## Walkthrough

The following provides an example of some of the features provided by the plugin.

Assuming all prerequisites have been met, first list all repositories and their associated charts using the `oc helm index` command:

```shell
oc helm index

REPOSITORY              NAME                            LATEST VERSION
redhat-helm-repo        ibm-b2bi-prod                   2.0.0
redhat-helm-repo        ibm-cpq-prod                    4.0.1
redhat-helm-repo        ibm-mongodb-enterprise-helm     0.1.0
redhat-helm-repo        ibm-object-storage-plugin       2.0.7
redhat-helm-repo        ibm-oms-ent-prod                6.0.0
redhat-helm-repo        ibm-oms-pro-prod                6.0.0
redhat-helm-repo        ibm-operator-catalog-enablement 1.1.0
redhat-helm-repo        ibm-sfg-prod                    2.0.0
redhat-helm-repo        nodejs                          0.0.1
redhat-helm-repo        nodejs-ex-k                     0.2.1
redhat-helm-repo        quarkus                         0.0.3
```

Next, create a new project for this walkthrough called `oc-helm-test`

```shell
oc new-project oc-helm-test
```

Next, install the `quarkus` chart from the `redhat-helm` repository and provide `quarkus` as the release name:

```shell
oc helm install quarkus redhat-helm-repo/quarkus

NAME: quarkus
NAMESPACE: oc-helm-test
LAST DEPLOYED: Mon Apr 26 05:35:55 2021
STATUS: deployed
REVISION: 1
```

A new build will be started and in a few moments, the resulting container will be deployed.

By default, the Build will make use of the _jvm_ mode of Quarkus. Native compilation can be enabled by setting the `build.mode` value to `native`. Upgrade the chart to modify the build mode:

```shell
oc helm upgrade quarkus redhat-helm-repo/quarkus --set build.mode=native

NAME: quarkus
NAMESPACE: oc-helm-test
LAST DEPLOYED: Mon Apr 26 05:44:50 2021
STATUS: deployed
REVISION: 2
```

The _quarkus_ BuildConfig will now be updated with _native_ compilation enabled.

Revert the changes by rolling back to the prior revision

```shell
oc helm rollback quarkus

Rollback was a success! Happy Helming!
```

Review the history of the release

```shell
oc helm history quarkus

REVISION        UPDATED                         STATUS          CHART           APP VERSION     DESCRIPTION
1               Mon Apr 26 05:35:55 2021        superseded      quarkus-0.0.3                   Install complete
2               Mon Apr 26 05:44:50 2021        superseded      quarkus-0.0.3                   Upgrade complete
3               Mon Apr 26 05:48:40 2021        deployed        quarkus-0.0.3                   Rollback to 1
```

Finally, uninstall the chart

```shell
oc helm uninstall quarkus

release "quarkus" uninstalled
```

You can also veify a chart. This allows users to validate a Helm chart URL and provides a report summary with the number of checks passed or failed.
You need to specify the profile which chart verifier will use to run the test. If no profile is provided default is used. With respect to console we recommend using `developer-console` profile. The following checks are enabled in `developer-console` profile:

| Checks | Description |
|---|---|
| IsHelmV3 | Checks that the given `uri` points to a Helm v3 chart. |
| HasReadme | Checks that the Helm chart contains the `README.md` file. |
| HasKubeVersion | Checks that the `Chart.yaml` file of the Helm chart includes the `kubeVersion` field (v1.0) and is a valid semantic version (v1.1). |
| ContainsValuesSchema | Checks that the Helm chart contains a JSON schema file (`values.schema.json`) to validate the `values.yaml` file in the chart. |
|  HelmLint | Checks that the chart is well formed by running the `helm lint` command. |
| ContainsValues | Checks that the Helm chart contains the `values`[¹](https://github.com/redhat-certification/chart-verifier/blob/main/docs/helm-chart-checks.md#-for-more-information-on-the-values-file-see-values-and-best-practices-for-using-values) file. |

example:

```shell
oc helm verify --chart-url https://github.com/openshift-helm-charts/charts/releases/download/redhat-mysql-sed-0.1.0/redhat-mysql-sed-0.1.0.tgz  --values provider=developer-console
NUMBER OF CHECKS PASSED: 6
NUMBER OF CHECKS FAILED: 0
```

## Development

1. Clone the repository and navigate to the project directory:

```shell
git clone https://github.com/sabre1041/oc-helm
cd oc-helm
```

2. Build the plugin

```shell
make build
```

The binary will be placed in the `bin` folder

3. Install the binary to your path

```shell
make install
```

4. Confirm the installation of the plugin

```shell
oc helm

OpenShift Command Line tool to interact with Helm capabilities.

Usage:
  oc-helm [command]
...

```