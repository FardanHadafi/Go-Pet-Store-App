package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// UserController defines the handlers for user-related HTTP endpoints.
type UserController interface {
	Register(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	Login(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	Update(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	ChangePassword(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	FindById(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	FindAll(w http.ResponseWriter, r *http.Request, params httprouter.Params)
}