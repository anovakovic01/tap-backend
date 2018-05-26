package auth

import "errors"

var (
	// ErrInvalidToken indicates that received token is invalid.
	ErrInvalidToken = errors.New("received invalid token")

	// ErrFailedClaimsExtraction indicates that auth service failed to extract claims
	// from received token.
	ErrFailedClaimsExtraction = errors.New("failed to extract claims from JWT")
)

// Service is used for communication with OAuth service.
type Service interface {
	// Authorizes token that client sent.
	Authorize(string) (Claims, error)
}

// Claims contains data that is retreived through oauth token.
type Claims struct {
	Iss           string `json:"iss"`
	Sub           string `json:"sub"`
	Iat           int64  `json:"iat,string"`
	Exp           int64  `json:"exp,string"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified,string"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Locale        string `json:"locale"`
}
