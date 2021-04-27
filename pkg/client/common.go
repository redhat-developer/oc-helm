package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/redhat-cop/oc-helm/pkg/types"
	"sigs.k8s.io/yaml"
)

func do(httpClient *http.Client, req *http.Request, v interface{}, jsonResponse bool, helmClient bool) (*http.Response, *types.HelmClientError) {
	resp, err := httpClient.Do(req)
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

				if helmClient {
					var helmServerError types.HelmServerError
					err = json.NewDecoder(resp.Body).Decode(&helmServerError)

					if err == nil {
						helmClientError.HelmServerError = &helmServerError
					}
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
