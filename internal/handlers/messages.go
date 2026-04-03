package handlers

import (
	"hoodhire-chat/internal/repositories"
	"hoodhire-chat/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type MessageHandler struct {
	Repo *repositories.MessageRepo
}

func NewMessageHandler(repo *repositories.MessageRepo) *MessageHandler {
	return &MessageHandler{Repo: repo}
}

func getUserID(c *fiber.Ctx) uint {
	raw := c.Locals("userID")
	switch v := raw.(type) {
	case float64:
		return uint(v)
	case uint:
		return v
	default:
		return 0
	}
}

// GET /messages/:userID?page=1&limit=20
func (h *MessageHandler) GetConversation(c *fiber.Ctx) error {
	userID := getUserID(c)
	otherUserID, err := strconv.ParseUint(c.Params("userID"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid user id"})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	messages, err := h.Repo.GetConversation(userID, uint(otherUserID), page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{"messages": messages})
}

// PATCH /messages/:userID/read
func (h *MessageHandler) MarkAsRead(c *fiber.Ctx) error {
	receiverID := getUserID(c)
	senderID, err := strconv.ParseUint(c.Params("userID"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid user id"})
	}
	if err := h.Repo.MarkAsRead(uint(senderID), receiverID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{"message": "messages marked as read"})
}

// GET /messages/unread
func (h *MessageHandler) GetUnreadCount(c *fiber.Ctx) error {
	userID := getUserID(c)

	count, err := h.Repo.GetUnreadCount(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{"unread_count": count})
}

// GET /messages/conversations
func (h *MessageHandler) GetConversationList(c *fiber.Ctx) error {
	userID := getUserID(c)

	ids, err := h.Repo.GetConversationList(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{"conversations": ids})
}

// GET /messages/unread/breakdown
func (h *MessageHandler) GetUnreadBreakdown(c *fiber.Ctx) error {
    userID := getUserID(c)
    counts, err := h.Repo.GetUnreadCountPerUser(userID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.Status(200).JSON(fiber.Map{"unread": counts})
}

func (h *MessageHandler) UploadFile(c *fiber.Ctx) error {
    file, err := c.FormFile("file")
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "file required"})
    }
    
    src, err := file.Open()
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "failed to open file"})
    }
    defer src.Close()
    
    // upload to cloudinary
    url, err := utils.UploadFile(src, file.Filename)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "failed to upload file"})
    }
    
    return c.Status(200).JSON(fiber.Map{"url": url})
}