package model

type LoginRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"Password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type BlogPost struct {
	ID uint64 `json:"id"`
	Title string `json:"title"`
	TagLine string `json:"tagline"`
	Tags string `json:"tags"`
	Content string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type HotTake struct {
	ID uint64 `json:"id"`
	Title string `json:"title"`
	Details string `json:"details"`
}

type User struct {
	ID uint64 `json:"id"`
	UserName string `json:"user_name"`
	Name string `json:"name"`
	Role string `json:"role"`
	CreatedAt string `json:"created_at"`
	CreatedAtInternal int
	Password string
}