package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
)

type VerifyResponse struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
}

func AuthMiddleware(c *fiber.Ctx) error {
	// in fiber v2, get full authorization header
	token := c.Get("Authorization")
	log.Println("Authorization header length:", len(token))

	if token == "" {
		cookieToken := c.Cookies("access-token")
		log.Println("Cookie token length:", len(cookieToken))
		if cookieToken != "" {
			token = "Bearer " + cookieToken
		}
	}

	if token == "" || token == "Bearer" {
		log.Println("no valid token found")
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	req, err := http.NewRequest("GET", os.Getenv("MAIN_API_URL")+"/auth/verify", nil)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	req.Header.Set("Authorization", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("main API call failed:", err)
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	log.Println("main API response status:", resp.StatusCode)
	defer resp.Body.Close()

	var result VerifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("decode failed:", err)
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	log.Println("verified userID:", result.UserID)
	c.Locals("userID", float64(result.UserID))
	c.Locals("role", result.Role)
	return c.Next()
}
