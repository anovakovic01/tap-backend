package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/anovakovic01/tap-backend/auth"
	"github.com/anovakovic01/tap-backend/news"
	"github.com/go-zoo/bone"
)

var (
	errUnauthorized = errors.New("unauthorized access")
	authSvc         auth.Service
)

// UpdateHandler adds endpoints to HTTP handler.
func UpdateHandler(mux *bone.Mux, svc news.Service, as auth.Service) {
	authSvc = as

	mux.GetFunc("/news", all(svc))
	mux.GetFunc("/news/:id", one(svc))
}

func all(svc news.Service) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if err := authorize(req); err != nil {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		items := svc.All()

		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(items)
	}
}

func one(svc news.Service) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if err := authorize(req); err != nil {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		sid := bone.GetValue(req, "id")
		id, err := strconv.ParseInt(sid, 10, 64)
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		item, err := svc.One(id)
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(item)
	}
}

func authorize(req *http.Request) error {
	ahs := req.Header["Authorization"]
	if len(ahs) == 0 || ahs[0] == "" {
		return errUnauthorized
	}

	ah := ahs[0]

	if _, err := authSvc.Authorize(ah); err != nil {
		fmt.Println(err)
		return errUnauthorized
	}

	return nil
}
