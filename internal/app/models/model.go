// Package models defines common models and errors for the application.
package models

// URLRequest represents a request to shorten a URL.
type URLRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// URLResponse represents the response containing the shortened URL.
type URLResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// URL represents a mapping between a short URL and its original counterpart.
type URL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// CTXKey is the type used as a context key for storing user ID.
type CTXKey string

// All constants used in project
const (
	UserIDKey     CTXKey = "userID"         // UserIDKey is the specific key used in the context to store user ID.
	CertPEM       string = "cert.pem"       // CertPEM is the file name for TLS cert
	PrivateKeyPEM string = "privateKey.pem" // PrivateKeyPEM is the file name for TLS private key
)
