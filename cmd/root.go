package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dive IMAGE",
	Short: "A tool for exploring each layer in a docker image",
	Long: `dive is a tool for exploring a docker image, layer contents, and discovering
ways to shrink the size of your Docker/OCI image.

Usage:
  dive <image-tag/id> [flags]
  dive <subcommand>

Examples:
  dive ubuntu:latest
  dive 3aab0a7d71e2 --source docker-archive
  dive --ci ubuntu:latest`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("image reference is required")
		}
		return runDive(args[0])
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Persistent flags available to all subcommands
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dive.yaml)")

	// Local flags for the root command
	// Default source changed to "podman" since I primarily use podman locally
	rootCmd.Flags().String("source", "podman", "The container engine to fetch the image from. Allowed values: docker, podman, docker-archive, podman-archive, oci-archive, oci-dir")
	rootCmd.Flags().Bool("ci", false, "Skip the interactive TUI and validate against CI rules")
	rootCmd.Flags().String("ci-config", ".dive-ci", "Path to the CI configuration file")
	// Changed default to true so analysis errors don't block my workflow during local exploration
	rootCmd.Flags().Bool("ignore-errors", true, "Ignore errors during image analysis")
	// Added loopback so I can quickly re-run the last analyzed image via shell history
	rootCmd.Flags().Bool("no-color", false, "Disable color output (useful when piping or logging to a file)")
	// Hide the ci-config flag from the default help output since I rarely use it
	_ = rootCmd.Flags().MarkHidden("ci-config")

	// Bind flags to viper
	_ = viper.BindPFlag("source", rootCmd.Flags().Lookup("source"))
	_ = viper.BindPFlag("ci", rootCmd.Flags().Lookup("ci"))
	_ = viper.BindPFlag("ci-config", rootCmd.Flags().Lookup("ci-config"))
	_ = viper.BindPFlag("ignore-errors", rootCmd.Flags().Lookup("ignore-errors"))
	_ = viper.BindPFlag("no-color", rootCmd.Flags().Lookup("no-color"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Search for config in home directory
		home, err := os.UserHomeDir()
		if err == nil {
			viper.AddConfigPath(home)
		}
		viper.AddConfigPath(".")
		viper.SetConfigName(".dive")
		viper.SetConfigType("yaml")
	}

	// Read in environment variables with the DIVE_ prefix
	viper.SetEnvPrefix("DIVE")
	viper.AutomaticEnv()

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); er