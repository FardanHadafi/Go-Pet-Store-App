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
	petRepository := repository.NewPetRepository()
	petService := service.NewPetService(petRepository, db, validate)
	PetController := controller.NewPetController(petService)

	router := httprouter.New()

	router.GET("/api/pets", PetController.FindAll)
	router.POST("/api/pets", PetController.Create)
	router.GET("/api/pets/:petId", PetController.FindById)
	router.PUT("/api/pets/:petId", PetController.Update)
	router.DELETE("/api/pets/:petId", PetController.Delete)

	router.PanicHandler = exception.ErrorHandler

	server := http.Server{
		Addr: "localhost:3000",
		Handler: router,
	}

	err := server.ListenAndServe()
	helper.PanicIfError(err)
}