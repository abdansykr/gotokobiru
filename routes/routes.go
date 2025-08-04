package routes

import (
	"time"
	"tokobiru/controllers"
	"tokobiru/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(router *gin.Engine, db *mongo.Client) {
	// TERAPKAN MIDDLEWARE CORS DI SINI
	// Ini harus menjadi salah satu middleware pertama yang diterapkan.
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Mengizinkan semua origin (untuk pengembangan)
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Inisialisasi semua controller
	authController := controllers.NewAuthController(db)
	productController := controllers.NewProductController(db)
	cartController := controllers.NewCartController(db)
	orderController := controllers.NewOrderController(db)
	adminController := controllers.NewAdminController(db)
	userController := controllers.NewUserController(db)
	chatController := controllers.NewChatController(db)
	api := router.Group("/api/v1")
	{
		// RUTE BARU UNTUK CHATBOT
		chatbot := api.Group("/chatbot")
		{
			chatbot.POST("/ask", chatController.HandleChat)
		}

		// Rute untuk autentikasi (Register & Login)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}

		// Rute untuk produk
		products := api.Group("/products")
		{
			products.GET("", productController.GetProducts)
			products.GET("/:id", productController.GetProductByID)
			products.POST("", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), productController.CreateProduct)
			products.PUT("/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), productController.UpdateProduct)
			products.DELETE("/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), productController.DeleteProduct)
		}

		// Rute untuk keranjang belanja (hanya untuk customer)
		cart := api.Group("/cart", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("customer"))
		{
			cart.GET("", cartController.GetCart)
			cart.POST("", cartController.AddItemToCart)
			cart.PUT("", cartController.UpdateCartItem)
			cart.DELETE("/:productId", cartController.RemoveItemFromCart)
		}

		// Rute untuk pemesanan/order (hanya untuk customer)
		orders := api.Group("/orders", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("customer"))
		{
			orders.POST("/checkout", orderController.Checkout)
			orders.GET("", orderController.GetUserOrders)
			orders.GET("/:id", orderController.GetOrderByID)
		}

		// Rute khusus untuk dashboard admin
		admin := api.Group("/admin", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"))
		{
			admin.GET("/users", adminController.GetAllUsers)
			admin.GET("/orders", adminController.GetAllOrders)
			admin.PATCH("/orders/:id", adminController.UpdateOrderStatus)
		}

		// Rute untuk manajemen user (profil sendiri)
		user := api.Group("/user", middlewares.AuthMiddleware())
		{
			user.PUT("/profile", userController.UpdateUserProfile)
		}
	}
}
