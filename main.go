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
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
