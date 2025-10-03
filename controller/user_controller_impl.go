package controller

import (
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

func (uc *UserControllerImpl) Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req web.UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		uc.writeErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := uc.validator.Struct(req); err != nil {
		uc.writeErrorResponse(w, fmt.Sprintf("Validation failed: %v", err), http.StatusBadRequest)
		return
	}
	resp, err := uc.userService.Register(r.Context(), req)
	if err != nil {
		uc.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	uc.writeJSONResponse(w, resp, http.StatusCreated)
}

func (uc *UserControllerImpl) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req web.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		uc.writeErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := uc.validator.Struct(req); err != nil {
		uc.writeErrorResponse(w, fmt.Sprintf("Validation failed: %v", err), http.StatusBadRequest)
		return
	}
	resp, err := uc.userService.Login(r.Context(), req)
	if err != nil {
		uc.writeErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}
	uc.writeJSONResponse(w, resp, http.StatusOK)
}

func (uc *UserControllerImpl) Update(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
    authenticatedUserID, ok := middleware.GetUserIDFromContext(r.Context())
    if !ok {
        uc.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    targetUserID, err := strconv.Atoi(params.ByName("id"))
    if err != nil {
        uc.writeErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    // Only allow self-update
    if authenticatedUserID != targetUserID {
        uc.writeErrorResponse(w, "Forbidden", http.StatusForbidden)
        return
    }

    var req web.UserUpdateRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        uc.writeErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    if err := uc.validator.Struct(req); err != nil {
        uc.writeErrorResponse(w, fmt.Sprintf("Validation failed: %v", err), http.StatusBadRequest)
        return
    }

    // Call service with both ID and request
    resp, err := uc.userService.Update(r.Context(), targetUserID, req)
    if err != nil {
        uc.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
        return
    }

    uc.writeJSONResponse(w, resp, http.StatusOK)
}


func (uc *UserControllerImpl) ChangePassword(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
    authenticatedUserID, ok := middleware.GetUserIDFromContext(r.Context())
    if !ok {
        uc.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    targetUserID, err := strconv.Atoi(params.ByName("id"))
    if err != nil {
        uc.writeErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    if authenticatedUserID != targetUserID {
        uc.writeErrorResponse(w, "Forbidden", http.StatusForbidden)
        return
    }

    var req web.UserChangePasswordRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        uc.writeErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // inject the user ID from params
    req.Id = targetUserID

    if err := uc.validator.Struct(req); err != nil {
        uc.writeErrorResponse(w, fmt.Sprintf("Validation failed: %v", err), http.StatusBadRequest)
        return
    }

    if err := uc.userService.ChangePassword(r.Context(), req); err != nil {
        uc.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
        return
    }

    uc.writeJSONResponse(w, map[string]string{"message": "Password updated successfully"}, http.StatusOK)
}

func (uc *UserControllerImpl) Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
    authenticatedUserID, ok := middleware.GetUserIDFromContext(r.Context())
    if !ok {
        uc.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    targetUserID, err := strconv.Atoi(params.ByName("id"))
    if err != nil {
        uc.writeErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    // Only allow self-delete (unless admin logic is added later)
    if authenticatedUserID != targetUserID {
        uc.writeErrorResponse(w, "Forbidden", http.StatusForbidden)
        return
    }

    if err := uc.userService.Delete(r.Context(), targetUserID); err != nil {
        uc.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
        return
    }

    uc.writeJSONResponse(w, map[string]string{"message": "User deleted successfully"}, http.StatusOK)
}

func (uc *UserControllerImpl) FindById(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userID, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		uc.writeErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	resp, err := uc.userService.FindById(r.Context(), userID)
	if err != nil {
		uc.writeErrorResponse(w, err.Error(), http.StatusNotFound)
		return
	}
	uc.writeJSONResponse(w, resp, http.StatusOK)
}

func (uc *UserControllerImpl) FindAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	resp, err := uc.userService.FindAll(r.Context())
	if err != nil {
		uc.writeErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	uc.writeJSONResponse(w, resp, http.StatusOK)
}

func (uc *UserControllerImpl) RefreshToken(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// obtains token from Authorization header
	auth := r.Header.Get("Authorization")
	if auth == "" {
		uc.writeErrorResponse(w, "Authorization header required", http.StatusUnauthorized)
		return
	}
	parts := strings.Fields(auth)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		uc.writeErrorResponse(w, "Invalid authorization header", http.StatusUnauthorized)
		return
	}
	oldToken := parts[1]
	resp, err := uc.userService.RefreshToken(r.Context(), oldToken)
	if err != nil {
		uc.writeErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}
	uc.writeJSONResponse(w, resp, http.StatusOK)
}

// helpers
func (uc *UserControllerImpl) writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(web.WebResponse{
		Code:   statusCode,
		Status: http.StatusText(statusCode),
		Data:   data,
	})
}

func (uc *UserControllerImpl) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(web.WebResponse{
		Code:   statusCode,
		Status: http.StatusText(statusCode),
		Data:   map[string]string{"error": message},
	})
}
