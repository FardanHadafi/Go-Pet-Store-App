package main

import (
	"Go-PetStoreApp/app"
	"Go-PetStoreApp/controller"
	"Go-PetStoreApp/exception"
	"Go-PetStoreApp/helper"
	"Go-PetStoreApp/middleware"
	"Go-PetStoreApp/repository"
	"Go-PetStoreApp/service"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/julienschmidt/httprouter"
)

func main() {
	cfg := app.LoadConfig()
	db := app.NewDB(cfg)
	validate := validator.New()

	// Repositories
	userRepo := repository.NewUserRepository()
	petRepo := repository.NewPetRepository()

	// Services (user needs token expiry)
	userService := service.NewUserService(userRepo, db, validate, cfg.TokenExpiry)
	petService := service.NewPetService(petRepo, db)

	// Controllers
	userController := controller.NewUserController(userService)
	petController := controller.NewPetController(petService)

	// Middleware
	jwtMiddleware := middleware.NewJWTMiddleware()
	router := httprouter.New()

	// --- User endpoints ---
	router.POST("/api/users/register", userController.Register)
	router.POST("/api/users/login", userController.Login)
	router.POST("/api/auth/refresh", userController.RefreshToken)

	router.GET("/api/users/:id", jwtMiddleware.Authenticate(userController.FindById))
	router.PUT("/api/users/:id", jwtMiddleware.Authenticate(userController.Update))
	router.PATCH("/api/users/:id/password", jwtMiddleware.Authenticate(userController.ChangePassword))
	router.DELETE("/api/users/:id", jwtMiddleware.Authenticate(userController.Delete))
	router.GET("/api/users", jwtMiddleware.Authenticate(userController.FindAll))

	// Admin-only users
	router.GET("/api/admin/users", jwtMiddleware.Authenticate(jwtMiddleware.RequireRole("admin", userController.FindAll)))

	// --- Pet endpoints ---
	router.GET("/api/pets", jwtMiddleware.Authenticate(petController.FindAll))
	router.POST("/api/pets", jwtMiddleware.Authenticate(petController.Create))
	router.GET("/api/pets/:petId", jwtMiddleware.Authenticate(petController.FindById))
	router.PUT("/api/pets/:petId", jwtMiddleware.Authenticate(petController.Update))
	router.DELETE("/api/pets/:petId", jwtMiddleware.Authenticate(petController.Delete))

	// Admin-only pets
	router.GET("/api/admin/pets", jwtMiddleware.Authenticate(jwtMiddleware.RequireRole("admin", petController.FindAll)))

	// Panic handler for JSON error response
	router.PanicHandler = exception.ErrorHandler

	// Wrap with logging middleware
	httpHandler := middleware.CORS(middleware.LoggingMiddleware(router))

	server := http.Server{
		Addr:    "localhost:3000",
		Handler: httpHandler,
	}

	if err := server.ListenAndServe(); err != nil {
		helper.PanicIfError(err)
	}
}
