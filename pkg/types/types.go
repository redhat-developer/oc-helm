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

type HelmVerifierRequest struct {
	ChartUrl string                 `json:"chart_url"`
	Values   map[string]interface{} `json:"values"`
}

type ResultsReport struct {
	Passed   string   `json:"passed" yaml:"passed"`
	Failed   string   `json:"failed" yaml:"failed"`
	Messages []string `json:"message" yaml:"message"`
}

type ApiResult struct {
	VerifierApiResult ResultsReport `json:"results"`
}

type ReleaseSecret struct {
	SecretName string `json:"secret_name"`
}
