package cli

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(publishCmd)
}

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish your wallpaper",
	Long:  "Start publishing your current wallpaper so others can subscribe.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println("Publishing wallpaper...")
		return nil
	},
}
