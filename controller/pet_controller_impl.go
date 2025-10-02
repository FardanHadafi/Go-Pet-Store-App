package controller

import (
	"Go-PetStoreApp/helper"
	"Go-PetStoreApp/middleware"
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

func (c *PetControllerImpl) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    var request web.PetCreateRequest
    helper.ReadFromRequest(r, &request)

    userID, ok := middleware.GetUserIDFromContext(r.Context())
    if !ok {
        webResponse := web.WebResponse{Code: 401, Status: "Unauthorized"}
        helper.WriteToResponseBody(w, webResponse)
        return
    }

    petResponse := c.PetService.Create(r.Context(), request, userID)

    webResponse := web.WebResponse{Code: 201, Status: "Created", Data: petResponse}
    helper.WriteToResponseBody(w, webResponse)
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

func (c *PetControllerImpl) FindAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    userID, ok := middleware.GetUserIDFromContext(r.Context())
    if !ok {
        webResponse := web.WebResponse{Code: 401, Status: "Unauthorized"}
        helper.WriteToResponseBody(w, webResponse)
        return
    }

    pets := c.PetService.FindAll(r.Context(), userID)

    webResponse := web.WebResponse{Code: 200, Status: "OK", Data: pets}
    helper.WriteToResponseBody(w, webResponse)
}