package config

import "fmt"

const (
	AppName    = "oasishub-indexer"
	AppVersion = "0.8.3"
	GitCommit  = "-"
	GoVersion  = "1.14"
)

func VersionString() string {
	return fmt.Sprintf(
		"%s %s (git: %s, %s)",
		AppName,
		AppVersion,
		GitCommit,
		GoVersion,
	)
}
