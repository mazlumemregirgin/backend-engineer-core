package main

import (
	"week-01-layered-architecture/controller"
	"week-01-layered-architecture/repository"
	"week-01-layered-architecture/service"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	repo := repository.NewUserRepository()
	userService := service.NewUserService(repo)
	userController := controller.NewUserController(userService)

	app.Get("/users", userController.GetAll)
	app.Post("/users", userController.Create)

	app.Listen(":3000")
}
