package dto

type Tag struct {
	ID     *int   `json:"id"`
	Name   string `json:"name"`
	UserID *int   `json:"user_id"`
}

type User struct {
	ID   *int   `json:"id"`
	Name string `json:"name"`
	Tags []Tag  `json:"tags"`
}
