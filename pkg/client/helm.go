package client

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/redhat-cop/oc-helm/pkg/options"
	"github.com/redhat-cop/oc-helm/pkg/types"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
)

const (
	// #nosec G101
	OPENSHIFT_SESSION_TOKEN_NAME = "openshift-session-token"
	CSRF_TOKEN_NAME              = "csrf-token"
	CSRF_HEADER                  = "X-CSRFToken"
	ORIGIN_HEADER                = "Origin"
)

type HelmChartClient struct {
	consoleURL string
	httpClient *http.Client
	cookies    []*http.Cookie
	headers    map[string]string
	namespace  string
}

func NewHelmChartClient(commonOptions *options.CommandLineOption) (*HelmChartClient, error) {

	cookies := []*http.Cookie{}
	csrfToken := randomString(64)
	consoleURL := fmt.Sprintf("https://%s", commonOptions.ConsoleHostname)

	cookies = append(cookies, &http.Cookie{
		Name:   CSRF_TOKEN_NAME,
		Value:  csrfToken,
		Path:   "/",
		Secure: true,
		Domain: fmt.Sprintf(".%s", commonOptions.ConsoleHostname),
	})

	cookies = append(cookies, &http.Cookie{
		Name:   OPENSHIFT_SESSION_TOKEN_NAME,
		Value:  commonOptions.Token,
		Path:   "/",
		Secure: true,
		Domain: fmt.Sprintf(".%s", commonOptions.ConsoleHostname),
	})

	headers := map[string]string{
		"Content-Type": "application/json",
		CSRF_HEADER:    csrfToken,
		ORIGIN_HEADER:  consoleURL,
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			// #nosec G402
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	return &HelmChartClient{
		cookies:    cookies,
		headers:    headers,
		consoleURL: consoleURL,
		namespace:  commonOptions.Namespace,
		httpClient: httpClient,
	}, nil
}

func (c *HelmChartClient) newRequest(method string, contextPath string, body interface{}) (*http.Request, error) {
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

	// Add cookies to the request
	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}

	// Add Headers to the Request
	for key, value := range c.headers {
		req.Header.Add(key, value)
	}

	return req, nil

}

func (c *HelmChartClient) History(releaseName string) (*[]release.Release, error) {

	req, err := c.newRequest("GET", fmt.Sprintf("/api/helm/release/history?name=%s&ns=%s", releaseName, c.namespace), nil)

	if err != nil {
		return nil, err
	}

	var releaseHistory []release.Release

	_, helmClientError := do(c.httpClient, req, &releaseHistory, true, true)

	if helmClientError != nil {

		if helmClientError.HelmServerError != nil {
			return nil, fmt.Errorf("%s", helmClientError.HelmServerError.Error)
		}

		return nil, fmt.Errorf("Failed to list history. Release '%s': Status code: %d", releaseName, helmClientError.StatusCode)
	}

	return &releaseHistory, nil

}

func (c *HelmChartClient) CreateRelease(releaseName string, chartUrl string, values map[string]interface{}, upgrade bool) (*release.Release, error) {

	helmRequest := &types.HelmRequest{
		Name:      releaseName,
		Namespace: c.namespace,
		ChartUrl:  chartUrl,
		Values:    values,
	}

	var req *http.Request
	var err error

	if upgrade {
		req, err = c.newRequest("PUT", "/api/helm/release", helmRequest)
	} else {
		req, err = c.newRequest("POST", "/api/helm/release", helmRequest)
	}

	if err != nil {
		return nil, err
	}

	var release release.Release

	_, helmClientError := do(c.httpClient, req, &release, true, true)

	if helmClientError != nil {

		if helmClientError.HelmServerError != nil {
			return nil, fmt.Errorf("%s", helmClientError.HelmServerError.Error)
		}

		return nil, fmt.Errorf("Failed to create release '%s': Status code: %d", releaseName, helmClientError.StatusCode)
	}

	return &release, nil

}

func (c *HelmChartClient) CreateReleaseAsync(releaseName string, chartUrl string, values map[string]interface{}, upgrade bool) (*types.ReleaseSecret, error) {

	helmRequest := &types.HelmRequest{
		Name:      releaseName,
		Namespace: c.namespace,
		ChartUrl:  chartUrl,
		Values:    values,
	}

	var req *http.Request
	var err error

	if upgrade {
		req, err = c.newRequest("PUT", "/api/helm/release/async", helmRequest)
	} else {
		req, err = c.newRequest("POST", "/api/helm/release/async", helmRequest)
	}

	if err != nil {
		return nil, err
	}

	var release types.ReleaseSecret

	_, helmClientError := do(c.httpClient, req, &release, true, true)

	if helmClientError != nil {

		if helmClientError.HelmServerError != nil {
			return nil, fmt.Errorf("%s", helmClientError.HelmServerError.Error)
		}

		return nil, fmt.Errorf("Failed to create release '%s': Status code: %d", releaseName, helmClientError.StatusCode)
	}

	return &release, nil

}

func (c *HelmChartClient) Rollback(releaseName string, revision int) (*release.Release, error) {

	helmRequest := &types.HelmRequest{
		Name:      releaseName,
		Namespace: c.namespace,
		Version:   revision,
	}

	req, err := c.newRequest("PATCH", "/api/helm/release", helmRequest)

	if err != nil {
		return nil, err
	}

	var rollbackReleaseResponse release.Release

	_, helmClientError := do(c.httpClient, req, &rollbackReleaseResponse, true, true)

	if helmClientError != nil {

		if helmClientError.HelmServerError != nil {
			return nil, fmt.Errorf("%s", helmClientError.HelmServerError.Error)
		}

		return nil, fmt.Errorf("Failed to rollback release '%s': Status code: %d", releaseName, helmClientError.StatusCode)
	}

	return &rollbackReleaseResponse, nil

}

func (c *HelmChartClient) Uninstall(releaseName string) (*release.UninstallReleaseResponse, error) {

	req, err := c.newRequest("DELETE", fmt.Sprintf("/api/helm/release?name=%s&ns=%s", releaseName, c.namespace), nil)

	if err != nil {
		return nil, err
	}

	var uninstallReleaseResponse release.UninstallReleaseResponse

	_, helmClientError := do(c.httpClient, req, &uninstallReleaseResponse, true, true)

	if helmClientError != nil {

		if helmClientError.HelmServerError != nil {
			return nil, fmt.Errorf("%s", helmClientError.HelmServerError.Error)
		}

		return nil, fmt.Errorf("Failed to Uninstall release '%s': Status code: %d", releaseName, helmClientError.StatusCode)
	}

	return &uninstallReleaseResponse, nil

}

func (c *HelmChartClient) GetIndex() (*repo.IndexFile, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/api/helm/charts/index.yaml?namespace=%s", c.namespace), nil)

	if err != nil {
		return nil, err
	}

	var indexFile repo.IndexFile

	_, helmClientError := do(c.httpClient, req, &indexFile, false, true)
	if helmClientError != nil {
		return nil, helmClientError.Error
	}

	return &indexFile, nil

}

func (c *HelmChartClient) ListReleases(limitInfo bool) (*[]release.Release, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/api/helm/releases?ns=%s&limitInfo=%s", c.namespace, limitInfo), nil)

	if err != nil {
		return nil, err
	}

	var release []release.Release

	_, helmClientError := do(c.httpClient, req, &release, true, true)

	if helmClientError != nil {
		return nil, helmClientError.Error
	}

	return &release, nil

}

func (c *HelmChartClient) GetChart(url string) (*chart.Chart, error) {

	req, err := c.newRequest("GET", fmt.Sprintf("/api/helm/chart?url=%s", url), nil)

	if err != nil {
		return nil, err
	}

	var chart chart.Chart

	_, helmClientError := do(c.httpClient, req, &chart, true, true)

	if helmClientError != nil {
		return nil, helmClientError.Error
	}

	return &chart, nil

}

func (c *HelmChartClient) GetRelease(releaseName string) (*release.Release, error) {

	req, err := c.newRequest("GET", fmt.Sprintf("/api/helm/release?name=%s&ns=%s", releaseName, c.namespace), nil)

	if err != nil {
		return nil, err
	}

	var release release.Release

	_, helmClientError := do(c.httpClient, req, &release, true, true)

	if helmClientError != nil {

		if helmClientError.HelmServerError != nil {
			return nil, fmt.Errorf("%s", helmClientError.HelmServerError.Error)
		}

		return nil, fmt.Errorf("Failed to get release '%s': Status code: %d", releaseName, helmClientError.StatusCode)
	}

	return &release, nil

}

func (c *HelmChartClient) createPath(contextPath string) string {
	return fmt.Sprintf("%s%s", c.consoleURL, contextPath)
}

func (c *HelmChartClient) VerifyChart(chartUrl string, values map[string]interface{}) (*types.ApiResult, error) {

	helmRequest := &types.HelmVerifierRequest{
		ChartUrl: chartUrl,
		Values:   values,
	}
	var req *http.Request
	var err error
	req, err = c.newRequest(http.MethodPost, "/api/helm/verify", helmRequest)
	if err != nil {
		return nil, err
	}

	var result types.ApiResult

	_, helmClientError := do(c.httpClient, req, &result, true, true)
	if helmClientError != nil {

		if helmClientError.HelmServerError != nil {
			return nil, fmt.Errorf("%s", helmClientError.HelmServerError.Error)
		}
		return nil, fmt.Errorf("Failed to verify chart '%s': Status code: %d", chartUrl, helmClientError.StatusCode)
	}

	return &result, nil

}

func randomString(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(fmt.Sprintf("FATAL ERROR: Unable to get random bytes for session token: %v", err))
	}
	return base64.StdEncoding.EncodeToString(bytes)
}
