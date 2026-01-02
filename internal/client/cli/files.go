package cli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.io/khosbilegt/wallstream/internal/client/api"
)

func init() {
	rootCmd.AddCommand(filesCmd)
	filesCmd.AddCommand(filesUploadCmd)
}

var filesCmd = &cobra.Command{
	Use:   "files",
	Short: "File management commands",
	Long:  "Upload and manage wallpaper files.",
}

var filesUploadCmd = &cobra.Command{
	Use:   "upload <file-path>",
	Short: "Upload a wallpaper file",
	Long:  "Upload a wallpaper file to the server.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		baseURL, _ := cmd.Flags().GetString("server")
		username, _ := cmd.Flags().GetString("username")
		apiKey, _ := cmd.Flags().GetString("api-key")

		if baseURL == "" {
			baseURL = "http://localhost:8080"
		}
		if username == "" || apiKey == "" {
			return fmt.Errorf("username and api-key are required for authenticated commands")
		}

		client := api.NewClient(baseURL, username, apiKey)
		ctx := context.Background()

		result, err := client.UploadWallpaper(ctx, filePath)
		if err != nil {
			return fmt.Errorf("failed to upload file: %w", err)
		}

		output, _ := json.MarshalIndent(result, "", "  ")
		cmd.Println(string(output))
		return nil
	},
}
