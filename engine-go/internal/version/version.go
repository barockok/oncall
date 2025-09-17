package version

var version = "dev-oss"

// Version returns the current build version. Can be overridden at build time using -ldflags.
func Version() string { return version }