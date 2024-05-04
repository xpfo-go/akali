package version

import "runtime"

var (
	Version   = "1.0.0"
	Commit    = "none"
	BuildTime = "unknown"
	GoVersion = runtime.Version()
)
