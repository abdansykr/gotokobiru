package main

import (
	"context"
	"log"
	"time"
	"tokobiru/config"
	"tokobiru/database"
	"tokobiru/models"
	"tokobiru/services"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	log.Println("Starting seeder...")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config for seeder: %v", err)
	}

	database.ConnectDB(cfg.MongoURI, cfg.MongoDatabase)

	seedUsers()
	seedProducts()

	log.Println("Seeding completed successfully!")
}

func seedUsers() {
	userCollection := database.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if admin user exists
	count, err := userCollection.CountDocuments(ctx, bson.M{"email": "admin@tokobiru.com"})
	if err != nil {
		log.Fatalf("Failed to check for admin user: %v", err)
	}

	if count > 0 {
		log.Println("Users already seeded. Skipping.")
		return
	}

	hashedPasswordAdmin, _ := services.HashPassword("admin123")
	hashedPasswordCustomer, _ := services.HashPassword("customer123")
	now := time.Now()

	users := []interface{}{
		models.User{
			ID:        primitive.NewObjectID(),
			Name:      "Admin User",
			Email:     "admin@tokobiru.com",
			Password:  hashedPasswordAdmin,
			Role:      "admin",
			CreatedAt: now,
			UpdatedAt: now,
		},
		models.User{
			ID:        primitive.NewObjectID(),
			Name:      "Customer Satu",
			Email:     "customer1@tokobiru.com",
			Password:  hashedPasswordCustomer,
			Role:      "customer",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	_, err = userCollection.InsertMany(ctx, users)
	if err != nil {
		log.Fatalf("Failed to seed users: %v", err)
	}
	log.Println("Users seeded.")
}

func seedProducts() {
	productCollection := database.GetCollection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := productCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Fatalf("Failed to check for products: %v", err)
	}

	if count > 0 {
		log.Println("Products already seeded. Skipping.")
		return
	}

	now := time.Now()
	products := []interface{}{
		models.Product{
			ID:          primitive.NewObjectID(),
			Name:        "Kaos Polos Biru Dongker",
			Description: "Kaos katun combed 30s, nyaman dan adem.",
			Price:       85000,
			Stock:       100,
			Category:    "Pakaian",
			ImageURL:    "https://placehold.co/600x400/1E3A8A/FFFFFF?text=Kaos+Biru",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		models.Product{
			ID:          primitive.NewObjectID(),
			Name:        "Kemeja Flanel Kotak-kotak",
			Description: "Kemeja flanel lengan panjang, cocok untuk gaya kasual.",
			Price:       175000,
			Stock:       50,
			Category:    "Pakaian",
			ImageURL:    "https://placehold.co/600x400/9CA3AF/FFFFFF?text=Kemeja+Flanel",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		models.Product{
			ID:          primitive.NewObjectID(),
			Name:        "Celana Jeans Slim Fit",
			Description: "Celana jeans dengan bahan stretch yang nyaman.",
			Price:       250000,
			Stock:       75,
			Category:    "Celana",
			ImageURL:    "https://placehold.co/600x400/374151/FFFFFF?text=Celana+Jeans",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		models.Product{
			ID:          primitive.NewObjectID(),
			Name:        "Topi Baseball Biru",
			Description: "Topi baseball dengan logo Toko Biru.",
			Price:       60000,
			Stock:       200,
			Category:    "Aksesoris",
			ImageURL:    "https://placehold.co/600x400/3B82F6/FFFFFF?text=Topi+Biru",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	_, err = productCollection.InsertMany(ctx, products)
	if err != nil {
		log.Fatalf("Failed to seed products: %v", err)
	}
	log.Println("Products seeded.")
}
