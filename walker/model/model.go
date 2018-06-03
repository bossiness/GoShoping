package model

// NoContent no content
type NoContent struct {
}

// PageRequest page request
type PageRequest struct {
	Offset int32 `json:"offset"`
	Limit  int32 `json:"limit"`
}

// IDRequest id request
type IDRequest struct {
	ID string `json:"_id"`
}

// UsernameRequest username
type UsernameRequest struct {
	Username string `json:"username"`
}