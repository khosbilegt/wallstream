package repository

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
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
