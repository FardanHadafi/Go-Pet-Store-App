package controller

import (
	"Go-PetStoreApp/helper"
	"Go-PetStoreApp/middleware"
	"Go-PetStoreApp/model/web"
	"Go-PetStoreApp/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator"
	"github.com/julienschmidt/httprouter"
)

type UserControllerImpl struct {
	userService service.UserService
	validator   *validator.Validate
}

func NewUserController(userService service.UserService) *UserControllerImpl {
	return &UserControllerImpl{
		userService: userService,
		validator:   validator.New(),
	}
}

// Register handles user registration (PUBLIC)
func (uc *UserControllerImpl) Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var request web.UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		uc.writeErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := uc.validator.Struct(request); err != nil {
		uc.writeErrorResponse(w, fmt.Sprintf("Validation failed: %v", err), http.StatusBadRequest)
		return
	}

	response, err := uc.userService.Register(r.Context(), request)
	if err != nil {
		if strings.Contains(err.Error(), "already registered") || 
		  strings.Contains(err.Error(), "already taken") {
			uc.writeErrorResponse(w, err.Error(), http.StatusConflict)
			return
		}
		uc.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	uc.writeJSONResponse(w, response, http.StatusCreated)
}

// Login handles user authentication (PUBLIC)
func (uc *UserControllerImpl) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var request web.UserLoginRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		uc.writeErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := uc.validator.Struct(request); err != nil {
		uc.writeErrorResponse(w, fmt.Sprintf("Validation failed: %v", err), http.StatusBadRequest)
		return
	}

	response, err := uc.userService.Login(r.Context(), request)
	if err != nil {
		uc.writeErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}

	uc.writeJSONResponse(w, response, http.StatusOK)
}

// Update handles user profile updates (PROTECTED - users can only update their own profile)
func (uc *UserControllerImpl) Update(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Get authenticated user ID from context
	authenticatedUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		uc.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get target user ID from URL
	targetUserID, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		uc.writeErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Check if user is trying to update their own profile
	if authenticatedUserID != targetUserID {
		uc.writeErrorResponse(w, "Forbidden: You can only update your own profile", http.StatusForbidden)
		return
	}

	var request web.UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		uc.writeErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	request.Id = targetUserID

	if err := uc.validator.Struct(request); err != nil {
		uc.writeErrorResponse(w, fmt.Sprintf("Validation failed: %v", err), http.StatusBadRequest)
		return
	}

	response, err := uc.userService.Update(r.Context(), request)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			uc.writeErrorResponse(w, err.Error(), http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "already taken") {
			uc.writeErrorResponse(w, err.Error(), http.StatusConflict)
			return
		}
		uc.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	uc.writeJSONResponse(w, response, http.StatusOK)
}

// ChangePassword handles password changes (PROTECTED - users can only change their own password)
func (uc *UserControllerImpl) ChangePassword(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Get authenticated user ID from context
	authenticatedUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		uc.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get target user ID from URL
	targetUserID, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		uc.writeErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Check if user is trying to change their own password
	if authenticatedUserID != targetUserID {
		uc.writeErrorResponse(w, "Forbidden: You can only change your own password", http.StatusForbidden)
		return
	}

	var request web.UserChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		uc.writeErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	request.Id = targetUserID

	if err := uc.validator.Struct(request); err != nil {
		uc.writeErrorResponse(w, fmt.Sprintf("Validation failed: %v", err), http.StatusBadRequest)
		return
	}

	if err := uc.userService.ChangePassword(r.Context(), request); err != nil {
		if strings.Contains(err.Error(), "not found") {
			uc.writeErrorResponse(w, err.Error(), http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "incorrect") {
			uc.writeErrorResponse(w, err.Error(), http.StatusUnauthorized)
			return
		}
		uc.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	uc.writeSuccessResponse(w, "Password changed successfully")
}

// Delete handles user deletion (PROTECTED - users can only delete their own account)
func (uc *UserControllerImpl) Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Get authenticated user ID from context
	authenticatedUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		uc.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get target user ID from URL
	targetUserID, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		uc.writeErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Check if user is trying to delete their own account
	if authenticatedUserID != targetUserID {
		uc.writeErrorResponse(w, "Forbidden: You can only delete your own account", http.StatusForbidden)
		return
	}

	if err := uc.userService.Delete(r.Context(), targetUserID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			uc.writeErrorResponse(w, err.Error(), http.StatusNotFound)
			return
		}
		uc.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// FindById retrieves a user by ID (PROTECTED)
func (uc *UserControllerImpl) FindById(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userID, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		uc.writeErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	response, err := uc.userService.FindById(r.Context(), userID)
	if err != nil {
		uc.writeErrorResponse(w, err.Error(), http.StatusNotFound)
		return
	}

	uc.writeJSONResponse(w, response, http.StatusOK)
}

// FindAll retrieves all users (PROTECTED)
func (uc *UserControllerImpl) FindAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	response := uc.userService.FindAll(r.Context())
	uc.writeJSONResponse(w, response, http.StatusOK)
}

func (uc *UserControllerImpl) writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(web.WebResponse{
        Code:   statusCode,
        Status: http.StatusText(statusCode),
        Data:   data,
    })
}

func (uc *UserControllerImpl) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(web.WebResponse{
        Code:   statusCode,
        Status: http.StatusText(statusCode),
        Data:   map[string]string{"error": message},
    })
}


func (uc *UserControllerImpl) writeSuccessResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": message,
	})
}

// RefreshToken generates a new JWT token using the current valid token
func (uc *UserControllerImpl) RefreshToken(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Get user info from context (already authenticated by middleware)
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		uc.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	email, ok := middleware.GetEmailFromContext(r.Context())
	if !ok {
		uc.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Generate new token
	newToken, err := helper.GenerateToken(userID, email, uc.userService.(*service.UserServiceImpl).TokenExpiry)
    if err != nil {
        uc.writeErrorResponse(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }

	response := map[string]string{
		"token": newToken,
		"message": "Token refreshed successfully",
	}

	uc.writeJSONResponse(w, response, http.StatusOK)
}