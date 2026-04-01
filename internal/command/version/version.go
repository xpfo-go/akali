package version

import "runtime"

var (
	Version   = "dev"
	Commit    = "none"
	BuildTime = "unknown"
	GoVersion = runtime.Version()
)
