package core

import (
	"runtime/debug"
)

var (
	version string
	commit  string
)

func runtimeVersion() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}

	for _, dep := range bi.Deps {
		if dep.Path == "github.com/gosthome/gosthome" {
			if dep.Version == "" {
				return "dev"
			}
			return dep.Version
		}
	}
	return "dev"
}

func Version() string {
	if version != "" {
		return version
	}
	return runtimeVersion()
}

func Commit() string {
	return commit
}
