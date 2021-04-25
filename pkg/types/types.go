package types

import (
	"helm.sh/helm/v3/pkg/repo"
)

type HelmClientError struct {
	StatusCode      int
	Message         string
	Error           error
	HelmServerError *HelmServerError
}

type HelmServerError struct {
	Error string `json:"error"`
}

type HelmRequest struct {
	Name      string                 `json:"name"`
	Namespace string                 `json:"namespace"`
	ChartUrl  string                 `json:"chart_url"`
	Values    map[string]interface{} `json:"values"`
	Version   int                    `json:"version"`
}

type ChartVersionRepository struct {
	Repository    string
	Chart         string
	ChartVersions repo.ChartVersions
}
