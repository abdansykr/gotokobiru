package controllers

import (
	"context"
	"net/http"
	"time"
	"tokobiru/database"
	"tokobiru/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserController untuk mengelola operasi terkait user
type UserController struct {
	db *mongo.Client
}

// NewUserController membuat instance baru dari UserController
func NewUserController(db *mongo.Client) *UserController {
	return &UserController{db: db}
}

// UpdateUserProfile mengizinkan pengguna yang sudah login untuk memperbarui
// nama atau password mereka sendiri.
func (uc *UserController) UpdateUserProfile(c *gin.Context) {
	// Mengambil userID dari token yang sudah divalidasi oleh middleware
	userIDHex, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}
	userID, _ := primitive.ObjectIDFromHex(userIDHex.(string))

	// Menangkap data dari body request
	var req struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userCollection := database.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Mempersiapkan field yang akan di-update
	updateFields := bson.M{}
	if req.Name != "" {
		updateFields["name"] = req.Name
	}

	if req.Password != "" {
		// Validasi sederhana untuk panjang password
		if len(req.Password) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 characters long"})
			return
		}
		// Hash password baru sebelum disimpan
		hashedPassword, err := services.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
			return
		}
		updateFields["password"] = hashedPassword
	}

	// Jika tidak ada data yang dikirim, kembalikan error
	if len(updateFields) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update provided"})
		return
	}

	// Selalu perbarui timestamp `updated_at`
	updateFields["updated_at"] = time.Now()

	updateQuery := bson.M{"$set": updateFields}

	// Menjalankan query update ke database
	result, err := userCollection.UpdateOne(ctx, bson.M{"_id": userID}, updateQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user profile"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User profile updated successfully"})
}
