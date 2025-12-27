package repository

type User struct {
	ID        string `json:"id" bson:"id"`
	Username  string `json:"username" bson:"username"`
	Password  string `json:"password" bson:"password"`
	APIKey    string `json:"api_key,omitempty" bson:"api_key"`
	CreatedAt int64  `json:"created_at" bson:"created_at"`
	UpdatedAt int64  `json:"updated_at" bson:"updated_at"`
}

type PublisherState struct {
	ID        string `json:"id"`
	Hash      string `json:"hash"`
	URL       string `json:"url"`
	Timestamp int64  `json:"timestamp"`
}

type SubscriberState struct {
	ID        string `json:"id"`
	Hash      string `json:"hash"`
	URL       string `json:"url"`
	Timestamp int64  `json:"timestamp"`
}
