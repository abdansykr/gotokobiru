package controllers

import (
	"context"
	"net/http"
	"tokobiru/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatController struct {
	chatService *services.ChatService
}

// Struct untuk menangkap input dari frontend
type ChatInput struct {
	Prompt string `json:"prompt" binding:"required"`
}

// NewChatController menginisialisasi controller dengan service yang dibutuhkan
func NewChatController(db *mongo.Client) *ChatController {
	// Inisialisasi klien Gemini dan ChatService
	geminiModel := services.SetupGeminiClient(context.Background())
	chatService := services.NewChatService(db, geminiModel)

	return &ChatController{
		chatService: chatService,
	}
}

// HandleChat menangani request dari frontend dan memanggil service
func (cc *ChatController) HandleChat(c *gin.Context) {
	var input ChatInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pesan tidak boleh kosong"})
		return
	}

	// Memanggil service untuk mendapatkan jawaban dari LLM
	response, err := cc.chatService.GenerateResponse(c.Request.Context(), input.Prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses permintaan Anda"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reply": response})
}
