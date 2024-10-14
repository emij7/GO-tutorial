package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Todo struct{
	// Define el campo ID como entero y el campo Title como una cadena, json:"id" significa q al representarse en json, sea como "id"
	ID int `json:"id"`
	Title string `json:"title"`
	Completed bool `json:"completed"`
	Body string `json:"body"`
}

func main() {
	fmt.Println("Hello FRESH")
	//inicializo la app
	app := fiber.New()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	PORT := os.Getenv("PORT")

	todos := []Todo{}

	//Uso el puerto 4000, lo importo desde .env
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"message": "Hello, World!"})
	})

	//Create todo
	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo  := &Todo{}
        if err := c.BodyParser(&todo); err!= nil {
            return c.Status(400).JSON(fiber.Map{"error": err.Error()})
        }
		if todo.ID != 0{
			return c.Status(400).JSON(fiber.Map{"error": "ID must not be set"})
		}
		if todo.Title == ""{
			return c.Status(400).JSON(fiber.Map{"error": "Title is required"})
		}

		todo.ID = len(todos) + 1
        todos = append(todos, *todo)

        return c.Status(201).JSON(todo)
	})

	//Update todo
	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		//Atoi es convertir string a entero
        id, err := strconv.Atoi(c.Params("id"))
        if err!= nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
        }

        for index, todo := range todos {
            if todo.ID == id {
                todos[index].Completed = !todos[index].Completed
                return c.Status(200).JSON(todos[index])
            }
        }

        return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
    })

	//Delete todo, después ver si lo cambio por soft delete
	app.Delete("/api/todos/:id", func(ctx *fiber.Ctx) error {
        id, err := strconv.Atoi(ctx.Params("id"))
        if err!= nil {
            return ctx.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
        }
		for index, todo := range todos {
			if todo.ID == id {
                todos = append(todos[:index], todos[index+1:]...)
                return ctx.Status(204).SendString("Todo deleted")
            }
        }
		return ctx.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	})
	//List todos
	app.Get("/api/todos", func(c *fiber.Ctx) error {
        return c.JSON(todos)
    })

    //List todo by id
    app.Get("/api/todos/:id", func(c *fiber.Ctx) error {
        id, err := strconv.Atoi(c.Params("id"))
        if err!= nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
        }

        for _, todo := range todos {
            if todo.ID == id {
                return c.JSON(todo)
            }
        }

        return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
    })

	log.Fatal(app.Listen(":"+PORT))
}