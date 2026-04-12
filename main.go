package main

import (
	"fmt"
	"os"
	"strings"

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
		fmt.Fprintln(os.Stderr, "error:", err.Error())
		// exit code 2 for usage/input errors, 1 for all other errors
		exitCode := 1
		if isUsageError(err) {
			exitCode = 2
		}
		os.Exit(exitCode)
	}
}

// isUsageError returns true if the error appears to be caused by incorrect
// user input or invalid CLI usage rather than an internal failure.
func isUsageError(err error) bool {
	if err == nil {
		return false
	}
	// heuristic: cobra usage errors typically contain "usage" in the message
	msg := strings.ToLower(err.Error())
	return strings.HasPrefix(msg, "usage")
}
