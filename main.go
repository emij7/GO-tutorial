package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
type Todo struct{
	// Define el campo ID como entero y el campo Title como una cadena, json:"id" significa q al representarse en json, sea como "id"
	// bson es para mongoDB como el campo completo se va a llamar id
	// primitive.ObjectID es un tipo de dato de mongoDB, cambié int por eso.
	ID primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	Title string `json:"title"`
	Completed bool `json:"completed"`
	Body string `json:"body"`
}

var collection *mongo.Collection

func main(){
	fmt.Println("Hello FRESH")
    err := godotenv.Load(".env")
    if err!= nil {
        log.Fatal("Error loading.env file")
    }

	MONGODB_URI := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(MONGODB_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	if err = client.Ping(context.Background(), nil); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	collection = client.Database("golang_db").Collection("todos")

	app := fiber.New()

	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodo)
	app.Patch("/api/todos/:id", updateTodo)
	app.Delete("/api/todos/:id", deleteTodo)

	port := os.Getenv("PORT")
	if port == "" {
        port = "4000"
    }

	log.Fatal(app.Listen("0.0.0.0:"+port))
}

func getTodos(c *fiber.Ctx) error {
	//bson.M{} es un filtro vacio, por lo que se traen todos los documentos
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	var todos []Todo

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo Todo
        err := cursor.Decode(&todo)
        if err!= nil {
            log.Fatal(err)
        }
        todos = append(todos, todo)
    }
	return c.JSON(todos)		
}

func createTodo(c *fiber.Ctx) error {
	todo := new(Todo)
	if err := c.BodyParser(todo); err != nil {
		return err
	}
	//Esto es igual a if todo.ID != 0
	if todo.ID != primitive.NilObjectID {
		return c.Status(400).JSON(fiber.Map{"error": "ID must not be set"})
	}
	if todo.Title == ""{
		return c.Status(400).JSON(fiber.Map{"error": "Title is required"})
	}
	insertResult, err := collection.InsertOne(context.Background(), todo)
	if err != nil {
		return err
	}

	todo.ID = insertResult.InsertedID.(primitive.ObjectID)

	return c.JSON(todo)
}

func updateTodo(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}
	//colección actualizar uno, argumentos -> contexto, que voy a actualizar, que voy a cambiar en el documento
	_ , err = collection.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": bson.M{"completed": true}})
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"message": "Todo updated successfully"})
}

func deleteTodo(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
    }
    _, err = collection.DeleteOne(context.Background(), bson.M{"_id": id})
    if err != nil {
        return err
    }
    return c.JSON(fiber.Map{"message": "Todo deleted successfully"})
}