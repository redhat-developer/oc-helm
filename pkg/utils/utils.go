package utils

import (
	"fmt"
	"strings"

	"github.com/redhat-cop/oc-helm/pkg/constants"
)

func SplitRepositoryIndexKey(key string) (string, string, error) {
	parts := strings.Split(key, constants.OPENSHIFT_HELM_INDEX_SEPERATOR)

	if len(parts) != 2 {
		return "", "", fmt.Errorf("'%s does not contain expected separator", key)
	}

	return parts[1], parts[0], nil
}

func CreateRepositoryIndexKey(repository string, chart string) string {

	return fmt.Sprintf("%s%s%s", chart, constants.OPENSHIFT_HELM_INDEX_SEPERATOR, repository)
}
