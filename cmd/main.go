package main

import (
	"hoodhire-chat/database"
	"hoodhire-chat/internal/handlers"
	"hoodhire-chat/internal/middleware"
	"hoodhire-chat/internal/repositories"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}

	db := database.Connect()

	msgRepo := repositories.NewMessageRepo(db)
	msgHandler := handlers.NewMessageHandler(msgRepo)

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))

	
	app.Get("/ws",
		middleware.AuthMiddleware,
		websocket.New(handlers.HandleWebSocket(msgRepo)),
	)

	api := app.Group("/messages", middleware.AuthMiddleware)
	api.Get("/conversations", msgHandler.GetConversationList)
	api.Get("/unread", msgHandler.GetUnreadCount)
	api.Get("/unread/breakdown", msgHandler.GetUnreadBreakdown) 
	api.Post("/upload", msgHandler.UploadFile)                  
	api.Get("/:userID", msgHandler.GetConversation)             
	api.Patch("/:userID/read", msgHandler.MarkAsRead)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("chat service running on :%s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
