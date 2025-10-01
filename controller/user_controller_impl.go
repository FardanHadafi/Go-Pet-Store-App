package controller

import (
	"Go-PetStoreApp/model/web"
	"Go-PetStoreApp/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type UserControllerImpl struct {
	userService service.UserService
}

// NewUserController creates a new UserController instance
func NewUserController(userService service.UserService) *UserControllerImpl {
	return &UserControllerImpl{
		userService: userService,
	}
}

// Register handles user registration
func (uc *UserControllerImpl) Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var request web.UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	response, err := uc.userService.Register(r.Context(), request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Login handles user authentication
func (uc *UserControllerImpl) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var request web.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	response, err := uc.userService.Login(r.Context(), request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Update handles user profile updates
func (uc *UserControllerImpl) Update(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userID, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var request web.UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Set the user ID from URL parameter
	request.Id = userID

	response, err := uc.userService.Update(r.Context(), request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ChangePassword handles password changes
func (uc *UserControllerImpl) ChangePassword(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userID, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var request web.UserChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Set the user ID from URL parameter
	request.Id = userID

	if err := uc.userService.ChangePassword(r.Context(), request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Delete handles user deletion
func (uc *UserControllerImpl) Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userID, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := uc.userService.Delete(r.Context(), userID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// FindById retrieves a user by ID
func (uc *UserControllerImpl) FindById(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userID, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	response, err := uc.userService.FindById(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// FindAll retrieves all users
func (uc *UserControllerImpl) FindAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	response := uc.userService.FindAll(r.Context())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}