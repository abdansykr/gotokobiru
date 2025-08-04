package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"tokobiru/database"
	"tokobiru/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderController struct {
	db *mongo.Client
}

func NewOrderController(db *mongo.Client) *OrderController {
	return &OrderController{db: db}
}

// Checkout (Versi Baru Tanpa Transaksi)
func (oc *OrderController) Checkout(c *gin.Context) {
	// ... (Fungsi ini tidak berubah) ...
	userIDHex, _ := c.Get("userID")
	userID, _ := primitive.ObjectIDFromHex(userIDHex.(string))

	cartCollection := database.GetCollection("carts")
	productCollection := database.GetCollection("products")
	orderCollection := database.GetCollection("orders")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var cart models.Cart
	err := cartCollection.FindOne(ctx, bson.M{"userId": userID}).Decode(&cart)
	if err != nil || len(cart.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Keranjang kosong atau tidak ditemukan"})
		return
	}

	var orderItems []models.OrderItem
	var total float64

	for _, item := range cart.Items {
		var product models.Product
		err := productCollection.FindOne(ctx, bson.M{"_id": item.ProductID}).Decode(&product)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Produk dengan ID %s tidak ditemukan", item.ProductID.Hex())})
			return
		}
		if product.Stock < item.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Stok untuk produk %s tidak mencukupi", product.Name)})
			return
		}

		newStock := product.Stock - item.Quantity
		update := bson.M{"$set": bson.M{"stock": newStock}}
		_, err = productCollection.UpdateOne(ctx, bson.M{"_id": item.ProductID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Gagal memperbarui stok untuk produk %s", product.Name)})
			return
		}

		orderItems = append(orderItems, models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})
		total += product.Price * float64(item.Quantity)
	}

	newOrder := models.Order{
		ID:        primitive.NewObjectID(),
		OrderID:   fmt.Sprintf("TB-%d", time.Now().UnixNano()),
		UserID:    userID,
		Items:     orderItems,
		Total:     total,
		Status:    "baru",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = orderCollection.InsertOne(ctx, newOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat pesanan"})
		return
	}

	_, err = cartCollection.DeleteOne(ctx, bson.M{"_id": cart.ID})
	if err != nil {
		log.Printf("Peringatan: Gagal menghapus keranjang untuk user %s setelah checkout", userID.Hex())
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Checkout berhasil", "order": newOrder})
}

// GetUserOrders retrieves all orders for the logged-in user
func (oc *OrderController) GetUserOrders(c *gin.Context) {
	userIDHex, _ := c.Get("userID")
	userID, _ := primitive.ObjectIDFromHex(userIDHex.(string))

	orderCollection := database.GetCollection("orders")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := orderCollection.Find(ctx, bson.M{"userId": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	var orders []models.Order
	if err = cursor.All(ctx, &orders); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode orders"})
		return
	}

	// --- PERBAIKAN DI SINI ---
	// Pastikan kita selalu mengembalikan array, bukan 'null' jika tidak ada pesanan.
	if orders == nil {
		orders = make([]models.Order, 0)
	}

	c.JSON(http.StatusOK, orders)
}

// GetOrderByID retrieves a single order by its ID for the logged-in user
func (oc *OrderController) GetOrderByID(c *gin.Context) {
	// ... (Fungsi ini tidak berubah) ...
	orderID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	userIDHex, _ := c.Get("userID")
	userID, _ := primitive.ObjectIDFromHex(userIDHex.(string))

	orderCollection := database.GetCollection("orders")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var order models.Order
	err = orderCollection.FindOne(ctx, bson.M{"_id": orderID, "userId": userID}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order details"})
		return
	}

	c.JSON(http.StatusOK, order)
}
