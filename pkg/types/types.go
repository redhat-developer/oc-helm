package types

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
