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
	// ✅ Load all config
	cfg := app.LoadConfig()

	db := app.NewDB(cfg)
	validate := validator.New()

	// ✅ JWT middleware always consistent
	jwtMiddleware := middleware.NewJWTMiddleware()

	// Pet Dependencies
	petRepository := repository.NewPetRepository()
	petService := service.NewPetService(petRepository, db, validate)
	petController := controller.NewPetController(petService)

	// User Dependencies (TokenExpiry passed here)
	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository, db, validate, cfg.TokenExpiry)
	userController := controller.NewUserController(userService)

	router := httprouter.New()

	// Pet endpoints (GET public, modifications protected)
	router.GET("/api/pets", petController.FindAll)
	router.POST("/api/pets", jwtMiddleware.Authenticate(petController.Create))
	router.GET("/api/pets/:petId", petController.FindById)
	router.PUT("/api/pets/:petId", jwtMiddleware.Authenticate(petController.Update))
	router.DELETE("/api/pets/:petId", jwtMiddleware.Authenticate(petController.Delete))

	// User endpoints
	router.POST("/api/users/register", userController.Register)
	router.POST("/api/users/login", userController.Login)
	router.PUT("/api/users/:id", jwtMiddleware.Authenticate(userController.Update))
	router.PATCH("/api/users/:id/password", jwtMiddleware.Authenticate(userController.ChangePassword))
	router.DELETE("/api/users/:id", jwtMiddleware.Authenticate(userController.Delete))
	router.GET("/api/users/:id", jwtMiddleware.Authenticate(userController.FindById))
	router.GET("/api/users", jwtMiddleware.Authenticate(userController.FindAll))
	router.POST("/api/auth/refresh", jwtMiddleware.Authenticate(userController.RefreshToken))

	router.PanicHandler = exception.ErrorHandler

	server := http.Server{
		Addr:    "localhost:3000",
		Handler: router,
	}

	err := server.ListenAndServe()
	helper.PanicIfError(err)
}
