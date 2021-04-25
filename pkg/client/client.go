package client

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/redhat-cop/oc-helm/pkg/options"
	"github.com/redhat-cop/oc-helm/pkg/types"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
	"sigs.k8s.io/yaml"
)

const (
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

	_, helmClientError := c.do(req, &releaseHistory, true)

	if helmClientError != nil {

		if helmClientError.HelmServerError != nil {
			return nil, fmt.Errorf("%s", helmClientError.HelmServerError.Error)
		}

		return nil, fmt.Errorf("Failed to list history. Release '%s': Status code: %d", releaseName, helmClientError.StatusCode)
	}

	return &releaseHistory, nil

}

func (c *HelmChartClient) CreateRelease(releaseName string, chartUrl string, values map[string]interface{}) (*release.Release, error) {

	helmRequest := &types.HelmRequest{
		Name:      releaseName,
		Namespace: c.namespace,
		ChartUrl:  chartUrl,
		Values:    values,
	}

	req, err := c.newRequest("POST", "/api/helm/release", helmRequest)

	if err != nil {
		return nil, err
	}

	var release release.Release

	_, helmClientError := c.do(req, &release, true)

	if helmClientError != nil {

		if helmClientError.HelmServerError != nil {
			return nil, fmt.Errorf("%s", helmClientError.HelmServerError.Error)
		}

		return nil, fmt.Errorf("Failed to install release '%s': Status code: %d", releaseName, helmClientError.StatusCode)
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

	_, helmClientError := c.do(req, &rollbackReleaseResponse, true)

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

	_, helmClientError := c.do(req, &uninstallReleaseResponse, true)

	if helmClientError != nil {

		if helmClientError.HelmServerError != nil {
			return nil, fmt.Errorf("%s", helmClientError.HelmServerError.Error)
		}

		return nil, fmt.Errorf("Failed to Uninstall release '%s': Status code: %d", releaseName, helmClientError.StatusCode)
	}

	return &uninstallReleaseResponse, nil

}

func (c *HelmChartClient) GetIndex() (*repo.IndexFile, error) {
	req, err := c.newRequest("GET", "/api/helm/charts/index.yaml", nil)

	if err != nil {
		return nil, err
	}

	var indexFile repo.IndexFile

	_, helmClientError := c.do(req, &indexFile, false)

	if helmClientError != nil {
		return nil, helmClientError.Error
	}

	return &indexFile, nil

}

func (c *HelmChartClient) ListReleases() (*[]release.Release, error) {

	req, err := c.newRequest("GET", fmt.Sprintf("/api/helm/releases?ns=%s", c.namespace), nil)

	if err != nil {
		return nil, err
	}

	var release []release.Release

	_, helmClientError := c.do(req, &release, true)

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

	_, helmClientError := c.do(req, &chart, true)

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

	_, helmClientError := c.do(req, &release, true)

	if helmClientError != nil {

		if helmClientError.HelmServerError != nil {
			return nil, fmt.Errorf("%s", helmClientError.HelmServerError.Error)
		}

		return nil, fmt.Errorf("Failed to get release '%s': Status code: %d", releaseName, helmClientError.StatusCode)
	}

	return &release, nil

}

func (c *HelmChartClient) do(req *http.Request, v interface{}, jsonResponse bool) (*http.Response, *types.HelmClientError) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &types.HelmClientError{Error: err}
	}

	defer resp.Body.Close()

	statusCode := resp.StatusCode

	helmClientError := &types.HelmClientError{
		StatusCode: statusCode,
		Message:    http.StatusText(statusCode),
	}

	if v != nil {

		if jsonResponse {

			if statusCode > 399 {
				var helmServerError types.HelmServerError
				err = json.NewDecoder(resp.Body).Decode(&helmServerError)

				if err == nil {
					helmClientError.HelmServerError = &helmServerError
				}

				return nil, helmClientError

			}

			err = json.NewDecoder(resp.Body).Decode(v)
			if err != nil {
				return nil, &types.HelmClientError{Error: err, StatusCode: statusCode}
			}

		} else {

			if resp.StatusCode > 399 {

				return nil, &types.HelmClientError{StatusCode: statusCode, Message: http.StatusText(statusCode), Error: fmt.Errorf("%s", http.StatusText(resp.StatusCode))}
			}

			bodyBytes, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				return nil, &types.HelmClientError{Error: err, StatusCode: statusCode}
			}

			err = yaml.Unmarshal(bodyBytes, v)

			if err != nil {
				return resp, &types.HelmClientError{Error: err, StatusCode: statusCode}
			}

		}
	}

	return resp, nil
}

func (c *HelmChartClient) createPath(contextPath string) string {
	return fmt.Sprintf("%s%s", c.consoleURL, contextPath)
}

func randomString(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(fmt.Sprintf("FATAL ERROR: Unable to get random bytes for session token: %v", err))
	}
	return base64.StdEncoding.EncodeToString(bytes)
}
