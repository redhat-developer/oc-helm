module github.com/redhat-cop/oc-helm

go 1.15

require (
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.3
	helm.sh/helm/v3 v3.5.0
	k8s.io/apimachinery v0.20.2 // indirect
	k8s.io/cli-runtime v0.20.1
	k8s.io/helm v2.17.0+incompatible
	rsc.io/letsencrypt v0.0.3 // indirect
	sigs.k8s.io/yaml v1.2.0

)

replace (
	github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309
	github.com/opencontainers/runc => github.com/opencontainers/runc v1.0.0-rc8.0.20190926150303-84373aaa560b
)
