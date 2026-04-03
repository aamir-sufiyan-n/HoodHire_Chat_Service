package repositories

import (
	"hoodhire-chat/models"

	"gorm.io/gorm"
)

type MessageRepo struct {
	DB *gorm.DB
}

func NewMessageRepo(db *gorm.DB) *MessageRepo {
	return &MessageRepo{DB: db}
}

// save a message to DB
func (r *MessageRepo) SaveMessage(msg *models.Message) error {
	return r.DB.Create(msg).Error
}

// get conversation between two users
func (r *MessageRepo) GetConversation(userA, userB uint, page, limit int) ([]models.Message, error) {
	var messages []models.Message
	offset := (page - 1) * limit
	err := r.DB.Where(
		"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		userA, userB, userB, userA,
	).Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error

	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
        messages[i], messages[j] = messages[j], messages[i]
    }
	return messages, err
}

// mark all messages from a sender as read
func (r *MessageRepo) MarkAsRead(senderID, receiverID uint) error {
	return r.DB.Model(&models.Message{}).
		Where("sender_id = ? AND receiver_id = ? AND is_read = ?", senderID, receiverID, false).
		Update("is_read", true).Error
}

// get unread message count for a user
func (r *MessageRepo) GetUnreadCount(userID uint) (int64, error) {
	var count int64
	err := r.DB.Model(&models.Message{}).
		Where("receiver_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

// get list of all users the current user has chatted with
func (r *MessageRepo) GetConversationList(userID uint) ([]uint, error) {
	var ids []uint
	err := r.DB.Model(&models.Message{}).
		Where("sender_id = ? OR receiver_id = ?", userID, userID).
		Select("DISTINCT CASE WHEN sender_id = ? THEN receiver_id ELSE sender_id END", userID).
		Scan(&ids).Error
	return ids, err
}

func (r *MessageRepo) GetUnreadCountPerUser(receiverID uint) (map[uint]int64, error) {
    var results []struct {
        SenderID uint
        Count    int64
    }
    err := r.DB.Model(&models.Message{}).
        Select("sender_id, COUNT(*) as count").
        Where("receiver_id = ? AND is_read = ?", receiverID, false).
        Group("sender_id").
        Scan(&results).Error
    
    if err != nil {
        return nil, err
    }
    
    counts := make(map[uint]int64)
    for _, r := range results {
        counts[r.SenderID] = r.Count
    }
    return counts, nil
}