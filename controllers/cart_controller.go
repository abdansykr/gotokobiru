package controllers

import (
	"context"
	"net/http"
	"time"
	"tokobiru/database"
	"tokobiru/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CartController struct {
	db *mongo.Client
}

func NewCartController(db *mongo.Client) *CartController {
	return &CartController{db: db}
}

// Get the user's shopping cart
func (cc *CartController) GetCart(c *gin.Context) {
	userIDHex, _ := c.Get("userID")
	userID, _ := primitive.ObjectIDFromHex(userIDHex.(string))

	cartCollection := database.GetCollection("carts")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var cart models.Cart
	err := cartCollection.FindOne(ctx, bson.M{"userId": userID}).Decode(&cart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusOK, gin.H{"items": []models.CartItem{}}) // Return empty cart
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart"})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// Add an item to the cart
func (cc *CartController) AddItemToCart(c *gin.Context) {
	var req struct {
		ProductID string `json:"productId" binding:"required"`
		Quantity  int    `json:"quantity" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDHex, _ := c.Get("userID")
	userID, _ := primitive.ObjectIDFromHex(userIDHex.(string))
	productID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Product ID"})
		return
	}

	// Check if product exists and has enough stock
	var product models.Product
	productCollection := database.GetCollection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = productCollection.FindOne(ctx, bson.M{"_id": productID, "stock": bson.M{"$gte": req.Quantity}}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found or insufficient stock"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify product"})
		return
	}

	cartCollection := database.GetCollection("carts")
	// Find user's cart
	var cart models.Cart
	err = cartCollection.FindOne(ctx, bson.M{"userId": userID}).Decode(&cart)

	if err == mongo.ErrNoDocuments {
		// Create a new cart if not exists
		newCart := models.Cart{
			ID:     primitive.NewObjectID(),
			UserID: userID,
			Items: []models.CartItem{{
				ProductID: productID,
				Quantity:  req.Quantity,
				Name:      product.Name,
				Price:     product.Price,
				ImageURL:  product.ImageURL,
			}},
		}
		_, err = cartCollection.InsertOne(ctx, newCart)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
			return
		}
		c.JSON(http.StatusCreated, newCart)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart"})
		return
	}

	// Update existing cart
	// Check if product is already in the cart
	itemIndex := -1
	for i, item := range cart.Items {
		if item.ProductID == productID {
			itemIndex = i
			break
		}
	}

	if itemIndex != -1 {
		// Product exists, update quantity
		cart.Items[itemIndex].Quantity += req.Quantity
	} else {
		// Product does not exist, add new item
		cart.Items = append(cart.Items, models.CartItem{
			ProductID: productID,
			Quantity:  req.Quantity,
			Name:      product.Name,
			Price:     product.Price,
			ImageURL:  product.ImageURL,
		})
	}

	// Check stock for updated quantity
	if itemIndex != -1 {
		err = productCollection.FindOne(ctx, bson.M{"_id": productID, "stock": bson.M{"$gte": cart.Items[itemIndex].Quantity}}).Decode(&models.Product{})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Insufficient stock for updated quantity"})
			return
		}
	}

	// Save the updated cart
	_, err = cartCollection.UpdateOne(ctx, bson.M{"_id": cart.ID}, bson.M{"$set": bson.M{"items": cart.Items}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart"})
		return
	}
	c.JSON(http.StatusOK, cart)
}

// Update quantity of an item in the cart
func (cc *CartController) UpdateCartItem(c *gin.Context) {
	var req struct {
		ProductID string `json:"productId" binding:"required"`
		Quantity  int    `json:"quantity" binding:"required,gte=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDHex, _ := c.Get("userID")
	userID, _ := primitive.ObjectIDFromHex(userIDHex.(string))
	productID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Product ID"})
		return
	}

	// Check stock
	if req.Quantity > 0 {
		var product models.Product
		productCollection := database.GetCollection("products")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err = productCollection.FindOne(ctx, bson.M{"_id": productID, "stock": bson.M{"$gte": req.Quantity}}).Decode(&product)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found or insufficient stock"})
			return
		}
	}

	cartCollection := database.GetCollection("carts")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var filter bson.M
	var update bson.M

	if req.Quantity > 0 {
		// Update quantity of a specific item in the array
		filter = bson.M{"userId": userID, "items.productId": productID}
		update = bson.M{"$set": bson.M{"items.$.quantity": req.Quantity}}
	} else {
		// If quantity is 0, remove the item from the cart
		filter = bson.M{"userId": userID}
		update = bson.M{"$pull": bson.M{"items": bson.M{"productId": productID}}}
	}

	result, err := cartCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart item"})
		return
	}
	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart or item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart updated successfully"})
}

// Remove an item from the cart
func (cc *CartController) RemoveItemFromCart(c *gin.Context) {
	productIDHex := c.Param("productId")
	productID, err := primitive.ObjectIDFromHex(productIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	userIDHex, _ := c.Get("userID")
	userID, _ := primitive.ObjectIDFromHex(userIDHex.(string))

	cartCollection := database.GetCollection("carts")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"userId": userID}
	update := bson.M{"$pull": bson.M{"items": bson.M{"productId": productID}}}

	result, err := cartCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove item from cart"})
		return
	}
	if result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found in cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart"})
}
