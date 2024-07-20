package handlers

import (
	"lemmego/api"
)

func Register(r *api.Router) {
	r.Post("/register", AuthStoreHandler)
}
