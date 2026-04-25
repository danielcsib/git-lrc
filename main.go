// git-lrc: A tool for managing Git repositories with LRC (Last Recent Commit) functionality.
// Fork of HexmosTech/git-lrc with extended features.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version is set at build time via ldflags
	Version = "dev"
	// Commit is the git commit hash set at build time
	Commit = "none"
	// BuildDate is the build date set at build time
	BuildDate = "unknown"
)

// rootCmd is the base command for git-lrc
var rootCmd = &cobra.Command{
	Use:   "git-lrc",
	Short: "git-lrc — manage and inspect recent Git commits",
	Long: `git-lrc is a CLI tool that extends Git with last-recent-commit (LRC)
workflows. It helps developers quickly inspect, compare, and act on
the most recent commits across one or more repositories.`,
	SilenceUsage:  true,
	SilenceErrors: true, // handle errors ourselves for cleaner output
}

// versionCmd prints build version information
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("git-lrc %s (commit: %s, built: %s)\n", Version, Commit, BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
