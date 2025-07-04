package goicinga

import "fmt"

var (
	// VersionMajor major version
	VersionMajor = 0
	// VersionMinor minor version
	VersionMinor = 4
	// VersionPatch patch level
	VersionPatch = 0
)

var (
	// BuildInfo contains the build timestamp
	BuildInfo = "development"
	// Version string
	Version = func() string {
		return fmt.Sprintf("%v.%v.%v (%v)", VersionMajor, VersionMinor, VersionPatch, BuildInfo)
	}
)
