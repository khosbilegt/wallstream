package cli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.io/khosbilegt/wallstream/internal/client/api"
)

func init() {
	rootCmd.AddCommand(devicesCmd)
	devicesCmd.AddCommand(devicesCreateCmd)
	devicesCmd.AddCommand(devicesListCmd)
	devicesCmd.AddCommand(devicesGetCmd)
	devicesCmd.AddCommand(devicesDeleteCmd)
	devicesCmd.AddCommand(devicesUploadURLCmd)
}

var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Device management commands",
	Long:  "Manage publisher devices for wallpaper sharing.",
}

var devicesCreateCmd = &cobra.Command{
	Use:   "create <device-id>",
	Short: "Create a new publisher device",
	Long:  "Register a new device for publishing wallpapers.",
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

		result, err := client.CreatePublisherDevice(ctx, deviceID)
		if err != nil {
			return fmt.Errorf("failed to create device: %w", err)
		}

		output, _ := json.MarshalIndent(result, "", "  ")
		cmd.Println(string(output))
		return nil
	},
}

var devicesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all publisher devices",
	Long:  "List all publisher devices for the authenticated user.",
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

		devices, err := client.GetPublisherDevices(ctx)
		if err != nil {
			return fmt.Errorf("failed to list devices: %w", err)
		}

		output, _ := json.MarshalIndent(devices, "", "  ")
		cmd.Println(string(output))
		return nil
	},
}

var devicesGetCmd = &cobra.Command{
	Use:   "get <device-id>",
	Short: "Get a publisher device by ID",
	Long:  "Get details of a specific publisher device.",
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

		device, err := client.GetPublisherDeviceByDeviceID(ctx, deviceID)
		if err != nil {
			return fmt.Errorf("failed to get device: %w", err)
		}

		output, _ := json.MarshalIndent(device, "", "  ")
		cmd.Println(string(output))
		return nil
	},
}

var devicesDeleteCmd = &cobra.Command{
	Use:   "delete <device-id>",
	Short: "Delete a publisher device",
	Long:  "Delete a publisher device by its ID.",
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

		if err := client.DeletePublisherDeviceByDeviceID(ctx, deviceID); err != nil {
			return fmt.Errorf("failed to delete device: %w", err)
		}

		cmd.Printf("Device %s deleted successfully\n", deviceID)
		return nil
	},
}

var devicesUploadURLCmd = &cobra.Command{
	Use:   "upload-url <device-id>",
	Short: "Get upload URL for a device",
	Long:  "Get a presigned upload URL for uploading wallpapers to a device.",
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

		result, err := client.GetUploadURL(ctx, deviceID)
		if err != nil {
			return fmt.Errorf("failed to get upload URL: %w", err)
		}

		output, _ := json.MarshalIndent(result, "", "  ")
		cmd.Println(string(output))
		return nil
	},
}
