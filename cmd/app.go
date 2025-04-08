package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aakritigkmit/payment-gateway/internal/handlers"
	"github.com/aakritigkmit/payment-gateway/internal/repository"
	"github.com/aakritigkmit/payment-gateway/internal/routes"
	"github.com/aakritigkmit/payment-gateway/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Load environment variables
func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error while reading .env")
	}
}

// MongoDB Connection
func mongoConnection() *mongo.Client {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGO_URI")).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		log.Fatal("MongoDB ping failed:", err)
	}
	fmt.Println("Connected to MongoDB!")
	return client
}

// Initialize services and routes
func initializeApp() (*mongo.Client, *chi.Mux) {
	mongoClient := mongoConnection()
	collection := mongoClient.Database(os.Getenv("MONGO_DBNAME")).Collection(os.Getenv("MONGO_COLLECTION_NAME"))

	userRepo := repository.NewUserRepo(collection)
	authRepo := repository.NewAuthRepo(collection)
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(authRepo, userRepo)

	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/api", func(subRouter chi.Router) {
		subRouter.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Server is healthy"))
		})
		subRouter.Mount("/users", routes.UserRoutes(userHandler))
		subRouter.Mount("/auth", routes.AuthRoutes(authHandler))
	})

	return mongoClient, r
}
