package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/redhat-cop/oc-helm/pkg/options"
)

type OpenShiftClient struct {
	httpClient *http.Client
	headers    map[string]string
	server     string
}

func NewOpenShiftClient(commonOptions *options.CommandLineOption) (*OpenShiftClient, error) {

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", commonOptions.Token),
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			// #nosec G402
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	return &OpenShiftClient{
		headers:    headers,
		server:     commonOptions.APIServer,
		httpClient: httpClient,
	}, nil
}

func (c *OpenShiftClient) newRequest(method string, contextPath string, body interface{}) (*http.Request, error) {
	url, err := url.Parse(c.createPath(contextPath))

	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url.String(), buf)

	// Add Headers to the Request
	for key, value := range c.headers {
		req.Header.Add(key, value)
	}

	return req, nil

}

func (c *OpenShiftClient) DiscoverConsoleURL() (map[string]interface{}, error) {
	req, err := c.newRequest("GET", "/api/v1/namespaces/openshift-config-managed/configmaps/console-public", nil)

	if err != nil {
		return nil, err
	}

	var configMap map[string]interface{}

	_, openshiftClientError := do(c.httpClient, req, &configMap, false, true)

	if openshiftClientError != nil {
		return nil, openshiftClientError.Error
	}

	return configMap, nil

}

func (c *OpenShiftClient) createPath(contextPath string) string {
	return fmt.Sprintf("%s%s", c.server, contextPath)
}
