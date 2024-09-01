package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Completed bool               `json:"completed"`
	Body      string             `json:"body"`
}

var collection *mongo.Collection

func main() {
	fmt.Println("hello world")

	if os.Getenv("ENV") != "production" {
		// Load the .env file if not in production
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file:", err)
		}
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(MONGODB_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MONGODB ATLAS")

	collection = client.Database("golang_db").Collection("todos")

	app := fiber.New()

	app.Post("/api/upload", uploadImage)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	if os.Getenv("ENV") == "production" {
		app.Static("/", "./client/dist")
	}

	log.Fatal(app.Listen("0.0.0.0:" + port))

}

func uploadImage(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Failed to receive image"})
	}

	// GCPのバケット名を指定
	bucketName := "your-bucket-name"

	// GCPに画像をアップロード
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create GCP client"})
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)
	object := bucket.Object(file.Filename)
	writer := object.NewWriter(ctx)
	defer writer.Close()

	fileContent, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to open image file"})
	}
	defer fileContent.Close()

	if _, err := io.Copy(writer, fileContent); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to upload image to GCP"})
	}

	// バケット名をPostgresのDBに保存
	db, err := sql.Open("postgres", "your-postgres-connection-string")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to connect to Postgres"})
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO buckets (name) VALUES ($1)", bucketName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save bucket name to Postgres"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Image uploaded successfully", "bucket": bucketName, "filename": file.Filename})
}
