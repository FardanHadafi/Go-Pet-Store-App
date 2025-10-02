package controller

import (
	"Go-PetStoreApp/helper"
	"Go-PetStoreApp/model/web"
	"Go-PetStoreApp/service"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type PetControllerImpl struct{
	PetService service.PetService
}

func NewPetController(petService service.PetService) PetController {
	return &PetControllerImpl{
		PetService: petService,
	}
}

func (p *PetControllerImpl) Create(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	petCreateRequest := web.PetCreateRequest{}
	helper.ReadFromRequest(r, &petCreateRequest)

	petResponse := p.PetService.Create(r.Context(), petCreateRequest)
	webRespose := web.WebResponse{
		Code: http.StatusCreated,
		Status: "Created",
		Data: petResponse,
	}

	w.WriteHeader(http.StatusCreated)
	helper.WriteToResponseBody(w, webRespose)
}

func (p *PetControllerImpl) Update(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	petUpdateRequest := web.PetUpdateRequest{}
	helper.ReadFromRequest(r, &petUpdateRequest)

	petId := params.ByName("petId")
	id, err := strconv.Atoi(petId)
	helper.PanicIfError(err)

	petUpdateRequest.Id = id

	petResponse := p.PetService.Update(r.Context(), petUpdateRequest)
	webRespose := web.WebResponse{
		Code: http.StatusOK,
		Status: "OK",
		Data: petResponse,
	}

	helper.WriteToResponseBody(w, webRespose)
}

func (p *PetControllerImpl) Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	petId := params.ByName("petId")
	id, err := strconv.Atoi(petId)
	helper.PanicIfError(err)

	p.PetService.Delete(r.Context(), id)
	webRespose := web.WebResponse{
		Code: 204,
		Status: strconv.Itoa(http.StatusNoContent),
	}

	helper.WriteToResponseBody(w, webRespose)
}

func (p *PetControllerImpl) FindById(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	petId := params.ByName("petId")
	id, err := strconv.Atoi(petId)
	helper.PanicIfError(err)

	petResponse := p.PetService.FindById(r.Context(), id)
	webRespose := web.WebResponse{
		Code: http.StatusOK,
		Status: "OK",
		Data: petResponse,
	}

	w.WriteHeader(http.StatusOK)
	helper.WriteToResponseBody(w, webRespose)
}

func (p *PetControllerImpl) FindAll(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	petResponses := p.PetService.FindAll(r.Context())
	webRespose := web.WebResponse{
		Code: http.StatusOK,
		Status: "OK",
		Data: petResponses,
	}

	w.WriteHeader(http.StatusOK)
	helper.WriteToResponseBody(w, webRespose)
}