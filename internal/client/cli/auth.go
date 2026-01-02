package cli

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(authLogoutCmd)
	authCmd.AddCommand(authStatusCmd)
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands",
	Long:  "Manage authentication and credentials for Wallstream.",
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to Wallstream",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: implement login
		cmd.Println("Logging in...")
		return nil
	},
}

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out of Wallstream",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println("Logged out.")
		return nil
	},
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication status",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println("Not authenticated.")
		return nil
	},
}
