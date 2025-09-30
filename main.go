package main

import (
	"Go-PetStoreApp/app"
	"Go-PetStoreApp/controller"
	"Go-PetStoreApp/repository"
	"Go-PetStoreApp/service"

	"github.com/go-playground/validator"
	"github.com/julienschmidt/httprouter"
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
}