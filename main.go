package main

import (
	"Go-PetStoreApp/app"
	"Go-PetStoreApp/controller"
	"Go-PetStoreApp/exception"
	"Go-PetStoreApp/helper"
	"Go-PetStoreApp/repository"
	"Go-PetStoreApp/service"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

func main() {
	db := app.NewDB()
	validate := validator.New()

	// Pet Dependencies
	petRepository := repository.NewPetRepository()
	petService := service.NewPetService(petRepository, db, validate)
	petController := controller.NewPetController(petService)

	// User Dependencies
	userRepository := repository.NewUserRepository()
  userService := service.NewUserService(userRepository, db, validate)
  userController := controller.NewUserController(userService)

	router := httprouter.New()

	// Pet endpoint
	router.GET("/api/pets", petController.FindAll)
	router.POST("/api/pets", petController.Create)
	router.GET("/api/pets/:petId", petController.FindById)
	router.PUT("/api/pets/:petId", petController.Update)
	router.DELETE("/api/pets/:petId", petController.Delete)

	// User endpoint
	router.POST("/api/users/register", userController.Register)
	router.POST("/api/users/login", userController.Login)
	router.PUT("/api/users/:id", userController.Update)
	router.PATCH("/api/users/:id/password", userController.ChangePassword)
	router.DELETE("/api/users/:id", userController.Delete)
	router.GET("/api/users/:id", userController.FindById)
	router.GET("/api/users", userController.FindAll)

	router.PanicHandler = exception.ErrorHandler

	server := http.Server{
		Addr: "localhost:3000",
		Handler: router,
	}

	err := server.ListenAndServe()
	helper.PanicIfError(err)
}