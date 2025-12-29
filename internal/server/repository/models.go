package repository

type User struct {
	ID        string `json:"id" bson:"id"`
	Username  string `json:"username" bson:"username"`
	APIKey    string `json:"api_key,omitempty" bson:"api_key"`
	CreatedAt int64  `json:"created_at" bson:"created_at"`
	UpdatedAt int64  `json:"updated_at" bson:"updated_at"`
}

type PublishedWallpaper struct {
	ID        string `json:"id" bson:"id"`
	UserID    string `json:"user_id" bson:"user_id"`
	DeviceID  string `json:"device_id" bson:"device_id"`
	Hash      string `json:"hash" bson:"hash"`
	URL       string `json:"url" bson:"url"`
	CreatedAt int64  `json:"created_at" bson:"created_at"`
	UpdatedAt int64  `json:"updated_at" bson:"updated_at"`
}

type PublisherDevice struct {
	ID          string `json:"id" bson:"id"`
	UserID      string `json:"user_id" bson:"user_id"`
	DeviceID    string `json:"device_id" bson:"device_id"`
	WallpaperId string `json:"wallpaper_id" bson:"wallpaper_id"`
	CreatedAt   int64  `json:"created_at" bson:"created_at"`
	UpdatedAt   int64  `json:"updated_at" bson:"updated_at"`
}
