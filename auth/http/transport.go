package http

import (
	"encoding/json"
	"net/http"

	"github.com/anovakovic01/tap-backend/auth"
	"github.com/go-zoo/bone"
)

// MakeHandler returns http handler with defined auth endpoints.
func MakeHandler(svc auth.Service) http.Handler {
	mux := bone.New()
	mux.PostFunc("/tokens/verify", verify(svc))
	return mux
}

func verify(svc auth.Service) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var tokenReq tokenRequest
		if err := json.NewDecoder(req.Body).Decode(&tokenReq); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		claims, err := svc.Authorize(tokenReq.Token)
		if err == auth.ErrInvalidToken {
			res.WriteHeader(http.StatusForbidden)
			return
		}
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		json.NewEncoder(res).Encode(claims)
	}
}

type tokenRequest struct {
	Token string `json:"token"`
}
