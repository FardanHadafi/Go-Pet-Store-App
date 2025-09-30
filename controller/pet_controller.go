package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type PetController interface {
	Create(w http.ResponseWriter, r *http.Response, params httprouter.Params)
	Update(w http.ResponseWriter, r *http.Response, params httprouter.Params)
	Delete(w http.ResponseWriter, r *http.Response, params httprouter.Params)
	FindById(w http.ResponseWriter, r *http.Response, params httprouter.Params)
	FindAll(w http.ResponseWriter, r *http.Response, params httprouter.Params)
}