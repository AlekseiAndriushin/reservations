package api

import (
	"github.com/AlexeyAndryushin/reservations/types"
	"github.com/gofiber/fiber/v2"
)

func HandleGetUsers(c *fiber.Ctx) error {
	user := types.User{
		FirstName: "Aleksus",
		LastName: "Alekseev",
	}
	return c.JSON(user)
}

func HandleGetUser(c *fiber.Ctx) error {
	return c.JSON("Alekseus")
}