package cli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.io/khosbilegt/wallstream/internal/client/api"
)

func init() {
	rootCmd.AddCommand(usersCmd)
	usersCmd.AddCommand(usersRegisterCmd)
}

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "User management commands",
	Long:  "Manage user accounts and authentication.",
}

var usersRegisterCmd = &cobra.Command{
	Use:   "register <username>",
	Short: "Register a new user",
	Long:  "Register a new user account and receive an API key.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		username := args[0]
		baseURL, _ := cmd.Flags().GetString("server")
		if baseURL == "" {
			baseURL = "http://localhost:8080"
		}

		client := api.NewClient(baseURL, "", "")
		ctx := context.Background()

		result, err := client.RegisterUser(ctx, username)
		if err != nil {
			return fmt.Errorf("failed to register user: %w", err)
		}

		output, _ := json.MarshalIndent(result, "", "  ")
		cmd.Println(string(output))
		cmd.Printf("\nSave your API key: %s\n", result.APIKey)
		return nil
	},
}
