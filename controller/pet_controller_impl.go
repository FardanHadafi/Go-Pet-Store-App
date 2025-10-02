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

type PetControllerImpl struct {
	PetService service.PetService
}

func NewPetController(s service.PetService) *PetControllerImpl {
	return &PetControllerImpl{PetService: s}
}

func (p *PetControllerImpl) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req web.PetCreateRequest
	if err := helper.ReadFromRequestBody(r, &req); err != nil {
		helper.WriteToResponseBody(w, web.WebResponse{Code: http.StatusBadRequest, Status: "Bad Request", Data: err.Error()})
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		helper.WriteToResponseBody(w, web.WebResponse{Code: http.StatusUnauthorized, Status: "Unauthorized"})
		return
	}

	petResp, err := p.PetService.Create(r.Context(), req, userID)
	if err != nil {
		helper.WriteToResponseBody(w, web.WebResponse{Code: http.StatusInternalServerError, Status: "Internal Server Error", Data: err.Error()})
		return
	}

	helper.WriteToResponseBody(w, web.WebResponse{Code: http.StatusCreated, Status: "Created", Data: petResp})
}

func (p *PetControllerImpl) FindAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))
	species := q.Get("species")
	ownerParam := q.Get("owner_id") // admin can pass owner_id to filter

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	role, _ := middleware.GetRoleFromContext(r.Context())
	// default: user returns only their pets
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok && role != "admin" {
		helper.WriteToResponseBody(w, web.WebResponse{Code: http.StatusUnauthorized, Status: "Unauthorized"})
		return
	}

	// Allow admin to pass owner_id to view specific user; if owner_id omitted and admin, set userID=0 to get all
	if role == "admin" {
		if ownerParam != "" {
			parsedOwner, _ := strconv.Atoi(ownerParam)
			userID = parsedOwner
		} else {
			// userID == 0 signals repository to ignore owner filter and return all
			userID = 0
		}
	}

	petsResp, total, err := p.PetService.FindAllByUser(r.Context(), userID, page, limit, species)
	if err != nil {
		helper.WriteToResponseBody(w, web.WebResponse{Code: http.StatusInternalServerError, Status: "Internal Server Error", Data: err.Error()})
		return
	}

	resp := map[string]interface{}{
		"items": petsResp,
		"page":  page,
		"limit": limit,
		"total": total,
	}
	helper.WriteToResponseBody(w, web.WebResponse{Code: http.StatusOK, Status: "OK", Data: resp})
}


func (p *PetControllerImpl) FindById(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	petID, _ := strconv.Atoi(params.ByName("petId"))
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		helper.WriteToResponseBody(w, web.WebResponse{Code: http.StatusUnauthorized, Status: "Unauthorized"})
		return
	}

	petResp, err := p.PetService.FindById(r.Context(), petID, userID)
	if err != nil {
		helper.WriteToResponseBody(w, web.WebResponse{Code: http.StatusNotFound, Status: "Not Found", Data: err.Error()})
		return
	}
	helper.WriteToResponseBody(w, web.WebResponse{Code: http.StatusOK, Status: "OK", Data: petResp})
}

func (p *PetControllerImpl) Update(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	petId, _ := strconv.Atoi(params.ByName("petId"))
	var req web.PetUpdateRequest
	if err := helper.ReadFromRequestBody(r, &req); err != nil {
		helper.WriteToResponseBody(w, web.WebResponse{Code: http.StatusBadRequest, Status: "Bad Request", Data: err.Error()})
		return
	}
	req.Id = petId

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		helper.WriteToResponseBody(w, web.WebResponse{Code: http.StatusUnauthorized, Status: "Unauthorized"})
		return
	}

	petResp, err := p.PetService.Update(r.Context(), req, userID)
	if err != nil {
		helper.WriteToResponseBody(w, web.WebResponse{Code: http.StatusForbidden, Status: "Forbidden", Data: err.Error()})
		return
	}

	helper.WriteToResponseBody(w, web.WebResponse{Code: http.StatusOK, Status: "OK", Data: petResp})
}

func (p *PetControllerImpl) Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	petId, _ := strconv.Atoi(params.ByName("petId"))
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		helper.WriteToResponseBody(w, web.WebResponse{Code: http.StatusUnauthorized, Status: "Unauthorized"})
		return
	}

	if err := p.PetService.Delete(r.Context(), petId, userID); err != nil {
		helper.WriteToResponseBody(w, web.WebResponse{Code: http.StatusForbidden, Status: "Forbidden", Data: err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
