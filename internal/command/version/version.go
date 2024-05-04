package version

import "runtime"

var (
	Version   = "1.0.1"
	Commit    = "none"
	BuildTime = "2024-05-04 18:39:00"
	GoVersion = runtime.Version()
)
