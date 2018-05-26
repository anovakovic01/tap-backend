package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/anovakovic01/tap-backend/auth"
	"github.com/anovakovic01/tap-backend/events"
	"github.com/go-zoo/bone"
)

var (
	errUnauthorized = errors.New("unauthorized access")
	authSvc         auth.Service
)

// UpdateHandler adds endpoints to HTTP handler.
func UpdateHandler(mux *bone.Mux, svc events.Service, as auth.Service) {
	authSvc = as

	mux.PostFunc("/events", create(svc))
	mux.GetFunc("/events", all(svc))
	mux.GetFunc("/events/:id", one(svc))
	mux.PatchFunc("/events", update(svc))
}

func create(svc events.Service) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		claims, err := authorize(req)
		if err != nil {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		var event events.Event
		if err := json.NewDecoder(req.Body).Decode(&event); err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		event.Owner = claims.Email
		id, err := svc.Create(event)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		event.ID = id

		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(event)
	}
}

func all(svc events.Service) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if _, err := authorize(req); err != nil {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		items := svc.All()

		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(items)
	}
}

func one(svc events.Service) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if _, err := authorize(req); err != nil {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		sid := bone.GetValue(req, "id")
		id, err := strconv.ParseInt(sid, 10, 64)
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		event, err := svc.One(id)
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(event)
	}
}

func update(svc events.Service) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		claims, err := authorize(req)
		if err != nil {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		var event events.Event
		if err := json.NewDecoder(req.Body).Decode(&event); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		event.Owner = claims.Email
		if err := svc.Update(event); err != nil {
			if err == events.ErrNotFound {
				res.WriteHeader(http.StatusNotFound)
				return
			}
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(event)
	}
}

func authorize(req *http.Request) (auth.Claims, error) {
	ahs := req.Header["Authorization"]
	if len(ahs) == 0 || ahs[0] == "" {
		return auth.Claims{}, errUnauthorized
	}

	ah := ahs[0]

	claims, err := authSvc.Authorize(ah)
	if err != nil {
		fmt.Println(err)
		return auth.Claims{}, errUnauthorized
	}

	return claims, nil
}
