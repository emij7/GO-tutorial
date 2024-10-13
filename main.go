package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	fmt.Println("Hello FRESH")
	//inicializo la app
	app := fiber.New()
	//Uso el puerto 4000
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"message": "Hello, World!"})
	})

	log.Fatal(app.Listen(":4000"))
}