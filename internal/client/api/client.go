package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type Client struct {
	baseURL    string
	username   string
	apiKey     string
	httpClient *http.Client
}

func NewClient(baseURL, username, apiKey string) *Client {
	return &Client{
		baseURL:    baseURL,
		username:   username,
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

func (c *Client) newRequest(
	ctx context.Context,
	method string,
	path string,
	body io.Reader,
	contentType string,
) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}

	// Basic Auth header
	auth := c.username + ":" + c.apiKey
	req.Header.Set(
		"Authorization",
		"Basic "+base64.StdEncoding.EncodeToString([]byte(auth)),
	)

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	} else {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	return req, nil
}

// User operations

type RegisterUserRequest struct {
	Username string `json:"username"`
}

type RegisterUserResponse struct {
	Username string `json:"username"`
	APIKey   string `json:"api_key"`
}

func (c *Client) RegisterUser(ctx context.Context, username string) (*RegisterUserResponse, error) {
	reqBody := RegisterUserRequest{Username: username}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/api/users/register", bytes.NewBuffer(jsonData), "application/json")
	if err != nil {
		return nil, err
	}

	// Remove auth for public endpoint
	req.Header.Del("Authorization")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("register failed: %s", errResp["error"])
	}

	var result RegisterUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// File operations

type UploadWallpaperResponse struct {
	Filename string `json:"filename"`
}

func (c *Client) UploadWallpaper(ctx context.Context, filePath string) (*UploadWallpaperResponse, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Extract just the filename from the path
	filename := filepath.Base(filePath)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(part, file); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/api/files/upload", &body, writer.FormDataContentType())
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("upload failed: %s", errResp["error"])
	}

	var result UploadWallpaperResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Publisher device operations

type CreatePublisherDeviceRequest struct {
	DeviceID string `json:"device_id"`
}

type CreatePublisherDeviceResponse struct {
	DeviceID string `json:"device_id"`
}

func (c *Client) CreatePublisherDevice(ctx context.Context, deviceID string) (*CreatePublisherDeviceResponse, error) {
	reqBody := CreatePublisherDeviceRequest{DeviceID: deviceID}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/api/publisher/devices", bytes.NewBuffer(jsonData), "application/json")
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("create device failed: %s", errResp["error"])
	}

	var result CreatePublisherDeviceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

type PublisherDevice struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	DeviceID  string `json:"device_id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func (c *Client) GetPublisherDevices(ctx context.Context) ([]PublisherDevice, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/api/publisher/devices", nil, "")
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("get devices failed: %s", errResp["error"])
	}

	var result []PublisherDevice
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) GetPublisherDeviceByDeviceID(ctx context.Context, deviceID string) (*PublisherDevice, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/publisher/devices/%s", deviceID), nil, "")
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("get device failed: %s", errResp["error"])
	}

	var result PublisherDevice
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) DeletePublisherDeviceByDeviceID(ctx context.Context, deviceID string) error {
	req, err := c.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/api/publisher/devices/%s", deviceID), nil, "")
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("delete device failed: %s", errResp["error"])
	}

	return nil
}

type GetUploadURLResponse struct {
	UploadURL string `json:"upload_url"`
}

func (c *Client) GetUploadURL(ctx context.Context, deviceID string) (*GetUploadURLResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/publisher/devices/%s/upload-url", deviceID), nil, "")
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("get upload URL failed: %s", errResp["error"])
	}

	var result GetUploadURLResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Wallpaper operations

type PublishUploadedWallpaperRequest struct {
	Filename string `json:"filename"`
	DeviceID string `json:"device_id"`
}

type PublishUploadedWallpaperResponse struct {
	Message string `json:"message"`
}

func (c *Client) PublishUploadedWallpaper(ctx context.Context, deviceID, filename string) (*PublishUploadedWallpaperResponse, error) {
	reqBody := PublishUploadedWallpaperRequest{
		Filename: filename,
		DeviceID: deviceID,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/api/publisher/wallpaper", bytes.NewBuffer(jsonData), "application/json")
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("publish wallpaper failed: %s", errResp["error"])
	}

	var result PublishUploadedWallpaperResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

type PublishedWallpaper struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	DeviceID  string `json:"device_id"`
	Hash      string `json:"hash"`
	URL       string `json:"url"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func (c *Client) GetPublishedWallpapers(ctx context.Context) ([]PublishedWallpaper, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/api/publisher/wallpaper", nil, "")
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("get wallpapers failed: %s", errResp["error"])
	}

	var result []PublishedWallpaper
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) GetPublishedWallpapersByDeviceID(ctx context.Context, deviceID string) ([]PublishedWallpaper, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/publisher/wallpaper/%s", deviceID), nil, "")
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("get wallpapers by device failed: %s", errResp["error"])
	}

	var result []PublishedWallpaper
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) DeletePublishedWallpaperByHash(ctx context.Context, hash string) error {
	req, err := c.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/api/publisher/wallpaper/%s", hash), nil, "")
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("delete wallpaper failed: %s", errResp["error"])
	}

	return nil
}

func (c *Client) ServeWallpaper(ctx context.Context, deviceID string) ([]byte, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/wallpaper/%s", deviceID), nil, "")
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("serve wallpaper failed: %s", errResp["error"])
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
