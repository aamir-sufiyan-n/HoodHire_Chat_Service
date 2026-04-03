package handlers

import (
	"encoding/json"
	"fmt"
	"hoodhire-chat/internal/repositories"
	"hoodhire-chat/internal/ws"
	"hoodhire-chat/models"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/websocket/v2"
)

type IncomingMessage struct {
	ReceiverID uint   `json:"receiver_id"`
	Content    string `json:"content"`
	Type       string `json:"type"` // "text", "image", "pdf"
}

type OutgoingMessage struct {
	ID         uint   `json:"id"`
	SenderID   uint   `json:"sender_id"`
	ReceiverID uint   `json:"receiver_id"`
	Content    string `json:"content"`
	Type       string `json:"type"`
	IsRead     bool   `json:"is_read"`
	CreatedAt  string `json:"created_at"`
}

func checkBond(seekerUserID, hirerUserID uint) bool {
	url := fmt.Sprintf("%s/bonds/check?seeker_user_id=%d&hirer_user_id=%d",
		os.Getenv("MAIN_API_URL"), seekerUserID, hirerUserID)

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		return false
	}
	defer resp.Body.Close()

	var result struct {
		Active bool `json:"active"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false
	}
	return result.Active
}
func HandleWebSocket(repo *repositories.MessageRepo) func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		userIDRaw := c.Locals("userID")
		if userIDRaw == nil {
			log.Println("no userID in locals")
			return
		}

		var userID uint
		switch v := userIDRaw.(type) {
		case float64:
			userID = uint(v)
		case uint:
			userID = v
		default:
			log.Println("unknown userID type")
			return
		}

		role, _ := c.Locals("role").(string)
		log.Printf("user %d connected as %s", userID, role)

		ws.GlobalHub.Register(userID, c)
		defer ws.GlobalHub.Unregister(userID)

		for {
			_, raw, err := c.ReadMessage()
			if err != nil {
				log.Printf("user %d disconnected: %v", userID, err)
				break
			}

			var incoming IncomingMessage
			if err := json.Unmarshal(raw, &incoming); err != nil {
				log.Printf("invalid message format from user %d", userID)
				continue
			}

			if incoming.ReceiverID == 0 || incoming.Content == "" {
				continue
			}

			var seekerUserID, hirerUserID uint
			if role == "seeker" {
				seekerUserID = userID
				hirerUserID = incoming.ReceiverID
			} else {
				hirerUserID = userID
				seekerUserID = incoming.ReceiverID
			}

			if !checkBond(seekerUserID, hirerUserID) {
				errMsg, _ := json.Marshal(map[string]string{"error": "no active bond with this user"})
				c.WriteMessage(1, errMsg)
				continue
			}

			msg := &models.Message{
				SenderID:   userID,
				ReceiverID: incoming.ReceiverID,
				Content:    incoming.Content,
				Type:       incoming.Type,
			}
			if err := repo.SaveMessage(msg); err != nil {
				log.Printf("failed to save message: %v", err)
				continue
			}

			outgoing, _ := json.Marshal(OutgoingMessage{
				ID:         msg.ID,
				SenderID:   userID,
				ReceiverID: incoming.ReceiverID,
				Content:    incoming.Content,
				Type:       incoming.Type,
				IsRead:     false,
				CreatedAt:  msg.CreatedAt.String(),
			})

			ws.GlobalHub.Send(incoming.ReceiverID, outgoing)
			ws.GlobalHub.Send(userID, outgoing)
		}
	}
}
