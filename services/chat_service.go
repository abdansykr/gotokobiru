package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"tokobiru/models"

	"github.com/google/generative-ai-go/genai"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/api/option"
)

// ChatService menangani semua logika yang berhubungan dengan AI
type ChatService struct {
	db          *mongo.Client
	geminiModel *genai.GenerativeModel
}

// NewChatService membuat instance baru dari ChatService
func NewChatService(db *mongo.Client, model *genai.GenerativeModel) *ChatService {
	return &ChatService{
		db:          db,
		geminiModel: model,
	}
}

// SetupGeminiClient adalah fungsi helper untuk inisialisasi klien Gemini
func SetupGeminiClient(ctx context.Context) *genai.GenerativeModel {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable not set.")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Failed to create Gemini client: %v", err)
	}

	model := client.GenerativeModel("gemini-1.5-flash")
	return model
}

// GenerateResponse mengirimkan prompt ke Gemini dan mengembalikan jawaban
func (s *ChatService) GenerateResponse(ctx context.Context, prompt string) (string, error) {
	// --- LOGIKA BARU: Ambil semua data produk dari DB ---
	productCollection := s.db.Database("tokobiruDB").Collection("products")
	cursor, err := productCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error fetching products for AI context: %v", err)
		// Tetap lanjutkan tanpa konteks produk jika ada error
	}

	var products []models.Product
	if err = cursor.All(ctx, &products); err != nil {
		log.Printf("Error decoding products for AI context: %v", err)
	}

	// Ubah data produk menjadi format JSON yang bisa dibaca AI
	productContextBytes, err := json.Marshal(products)
	productContext := ""
	if err == nil {
		productContext = string(productContextBytes)
	}
	// --- AKHIR LOGIKA BARU ---

	// Buat prompt yang kaya dengan konteks produk
	fullPrompt := fmt.Sprintf(`
Anda adalah asisten AI untuk toko online bernama 'Toko Biru'. 
Tugas Anda adalah menjawab pertanyaan pelanggan dengan ramah, membantu, dan informatif berdasarkan data produk yang tersedia.
Selalu jawab dalam Bahasa Indonesia.

Berikut adalah data produk kami dalam format JSON:
%s

Berdasarkan data di atas, jawab pertanyaan pelanggan berikut: "%s"
`, productContext, prompt)

	resp, err := s.geminiModel.GenerateContent(ctx, genai.Text(fullPrompt))
	if err != nil {
		return "", fmt.Errorf("gagal menghasilkan konten dari Gemini: %w", err)
	}

	var replyText strings.Builder
	if len(resp.Candidates) > 0 {
		for _, part := range resp.Candidates[0].Content.Parts {
			if txt, ok := part.(genai.Text); ok {
				replyText.WriteString(string(txt))
			}
		}
	}

	if replyText.Len() == 0 {
		return "Maaf, saya tidak bisa memberikan jawaban saat ini.", nil
	}

	return replyText.String(), nil
}
