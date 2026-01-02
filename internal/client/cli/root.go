package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wallstream",
	Short: "Wallstream syncs wallpapers between people",
	Long: `Wallstream lets you publish your wallpaper and subscribe to others.
It runs as a background service and works across Windows, macOS, and Linux.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Default behavior when no subcommand is provided
		return cmd.Help()
	},
}

// Execute is the CLI entrypoint
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
