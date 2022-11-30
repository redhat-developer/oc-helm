package utils

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/redhat-cop/oc-helm/pkg/options"
	"k8s.io/helm/pkg/strvals"
	"sigs.k8s.io/yaml"
)

func MergeValues(commandLineOptions *options.CommandLineOption) (map[string]interface{}, error) {
	base := map[string]interface{}{}

	// User specified a values files via -f/--values
	for _, filePath := range commandLineOptions.ValueFiles {
		currentMap := map[string]interface{}{}

		bytes, err := readFile(filePath)
		if err != nil {
			return nil, err
		}

		if err := yaml.Unmarshal(bytes, &currentMap); err != nil {
			return nil, errors.Wrapf(err, "failed to parse %s", filePath)
		}
		// Merge with the previous map
		base = mergeMaps(base, currentMap)
	}

	// User specified a value via --set
	for _, value := range commandLineOptions.Values {
		if err := strvals.ParseInto(value, base); err != nil {
			return nil, errors.Wrap(err, "failed parsing --set data")
		}
	}

	// User specified a value via --set-string
	for _, value := range commandLineOptions.StringValues {
		if err := strvals.ParseIntoString(value, base); err != nil {
			return nil, errors.Wrap(err, "failed parsing --set-string data")
		}
	}

	// User specified a value via --set-file
	for _, value := range commandLineOptions.FileValues {
		reader := func(rs []rune) (interface{}, error) {
			bytes, err := readFile(string(rs))
			return string(bytes), err
		}
		if err := strvals.ParseIntoFile(value, base, reader); err != nil {
			return nil, errors.Wrap(err, "failed parsing --set-file data")
		}
	}

	return base, nil
}

func mergeMaps(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if v, ok := v.(map[string]interface{}); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[string]interface{}); ok {
					out[k] = mergeMaps(bv, v)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}

// readFile load a file from stdin, the local directory, or a remote file with a url.
func readFile(filePath string) ([]byte, error) {
	if strings.TrimSpace(filePath) == "-" {
		return ioutil.ReadAll(os.Stdin)
	}
	// #nosec G304
	return ioutil.ReadFile(filePath)

}
