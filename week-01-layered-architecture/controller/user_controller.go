package controller

import (
	"week-01-layered-architecture/model"
	"week-01-layered-architecture/service"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	service *service.UserService
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{service}
}

func (c *UserController) GetAll(ctx *fiber.Ctx) error {
	users := c.service.GetAllUsers()
	return ctx.JSON(users)
}

func (c *UserController) Create(ctx *fiber.Ctx) error {
	var user model.User
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	newUser, err := c.service.CreateUser(user)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(newUser)
}
