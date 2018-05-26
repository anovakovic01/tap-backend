package google

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/anovakovic01/tap-backend/auth"
)

var _ auth.Service = (*service)(nil)

// Service contains user authentication and authorization logic.
type service struct {
	clientID string
}

// NewService returns new auth service instance.
func NewService(clientID string) auth.Service {
	return service{clientID}
}

func (svc service) Authorize(token string) (auth.Claims, error) {
	url := fmt.Sprintf("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=%s", token)
	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return auth.Claims{}, auth.ErrInvalidToken
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Println(res.StatusCode)
		return auth.Claims{}, auth.ErrInvalidToken
	}
	var claims auth.Claims
	if err := json.NewDecoder(res.Body).Decode(&claims); err != nil {
		fmt.Println(err)
		return auth.Claims{}, auth.ErrFailedClaimsExtraction
	}
	return claims, nil
}
