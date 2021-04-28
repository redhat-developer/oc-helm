package version

import (
	"fmt"
	"runtime"
)

var version, gitCommit, buildDate string
var platform = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)

// BuildInfo stores static build information about the binary.
type BuildInfo struct {
	BuildDate string
	Compiler  string
	GitCommit string
	GoVersion string
	Platform  string
	Version   string
}

// GetBuildInfo returns build information about the binary
func GetBuildInfo() *BuildInfo {
	// These vars are set via -ldflags settings during 'go build'
	return &BuildInfo{
		Version:   version,
		GitCommit: gitCommit,
		BuildDate: buildDate,
		GoVersion: runtime.Version(),
		Compiler:  runtime.Compiler,
		Platform:  platform,
	}
}
