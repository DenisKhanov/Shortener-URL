package models

type URLRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}
type URLResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
type URL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
type CTXKey string

const UserIDKey CTXKey = "userID"
