package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.io/khosbilegt/wallstream/internal/client/api"
)

func init() {
	rootCmd.AddCommand(wallpapersCmd)
	wallpapersCmd.AddCommand(wallpapersPublishCmd)
	wallpapersCmd.AddCommand(wallpapersListCmd)
	wallpapersCmd.AddCommand(wallpapersGetByDeviceCmd)
	wallpapersCmd.AddCommand(wallpapersDeleteCmd)
	wallpapersCmd.AddCommand(wallpapersServeCmd)
}

var wallpapersCmd = &cobra.Command{
	Use:   "wallpapers",
	Short: "Wallpaper management commands",
	Long:  "Manage published wallpapers.",
}

var wallpapersPublishCmd = &cobra.Command{
	Use:   "publish <device-id> <filename>",
	Short: "Publish an uploaded wallpaper",
	Long:  "Publish a previously uploaded wallpaper to a device.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		deviceID := args[0]
		filename := args[1]
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

		result, err := client.PublishUploadedWallpaper(ctx, deviceID, filename)
		if err != nil {
			return fmt.Errorf("failed to publish wallpaper: %w", err)
		}

		output, _ := json.MarshalIndent(result, "", "  ")
		cmd.Println(string(output))
		return nil
	},
}

var wallpapersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all published wallpapers",
	Long:  "List all published wallpapers for the authenticated user.",
	RunE: func(cmd *cobra.Command, args []string) error {
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

		wallpapers, err := client.GetPublishedWallpapers(ctx)
		if err != nil {
			return fmt.Errorf("failed to list wallpapers: %w", err)
		}

		output, _ := json.MarshalIndent(wallpapers, "", "  ")
		cmd.Println(string(output))
		return nil
	},
}

var wallpapersGetByDeviceCmd = &cobra.Command{
	Use:   "get-by-device <device-id>",
	Short: "Get wallpapers by device ID",
	Long:  "Get all published wallpapers for a specific device.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		deviceID := args[0]
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

		wallpapers, err := client.GetPublishedWallpapersByDeviceID(ctx, deviceID)
		if err != nil {
			return fmt.Errorf("failed to get wallpapers: %w", err)
		}

		output, _ := json.MarshalIndent(wallpapers, "", "  ")
		cmd.Println(string(output))
		return nil
	},
}

var wallpapersDeleteCmd = &cobra.Command{
	Use:   "delete <hash>",
	Short: "Delete a published wallpaper",
	Long:  "Delete a published wallpaper by its hash.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		hash := args[0]
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

		if err := client.DeletePublishedWallpaperByHash(ctx, hash); err != nil {
			return fmt.Errorf("failed to delete wallpaper: %w", err)
		}

		cmd.Printf("Wallpaper %s deleted successfully\n", hash)
		return nil
	},
}

var wallpapersServeCmd = &cobra.Command{
	Use:   "serve <device-id> [output-file]",
	Short: "Download/serve a wallpaper",
	Long:  "Download the latest wallpaper for a device. If output-file is provided, saves to file; otherwise prints to stdout.",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		deviceID := args[0]
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

		data, err := client.ServeWallpaper(ctx, deviceID)
		if err != nil {
			return fmt.Errorf("failed to serve wallpaper: %w", err)
		}

		if len(args) == 2 {
			// Save to file
			outputFile := args[1]
			if err := os.WriteFile(outputFile, data, 0644); err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}
			cmd.Printf("Wallpaper saved to %s\n", outputFile)
		} else {
			// Print to stdout
			os.Stdout.Write(data)
		}

		return nil
	},
}
