package main

import (
	"fmt"
	"os"

	"github.com/anchore/dive/cmd"
)

// version information set by ldflags at build time
var (
	version = "dev"
	buildDate = "not set"
	gitCommit = "not set"
	gitDescription = "not set"
)

func main() {
	if err := cmd.Execute(
		cmd.BuildInfo{
			Version:        version,
			BuildDate:      buildDate,
			GitCommit:      gitCommit,
			GitDescription: gitDescription,
		},
	); err != nil {
		// print the error to stderr and exit with a non-zero status code
		// using exit code 1 to indicate a general error (could be more specific in the future)
		fmt.Fprintln(os.Stderr, "error:", err.Error())
		// NOTE: consider mapping specific error types to distinct exit codes
		// e.g. exit code 2 for usage errors, 3 for image-not-found, etc.
		os.Exit(1)
	}
}
