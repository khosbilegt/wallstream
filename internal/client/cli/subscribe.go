package cli

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(subscribeCmd)
}

var subscribeCmd = &cobra.Command{
	Use:   "subscribe <username>",
	Short: "Subscribe to a user's wallpaper",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		username := args[0]
		cmd.Printf("Subscribing to %s...\n", username)
		return nil
	},
}
