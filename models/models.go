package models

import (
	"gorm.io/gorm"
)

type Message struct {
    gorm.Model
    SenderID   uint   `gorm:"index;not null" json:"sender_id"`
    ReceiverID uint   `gorm:"index;not null" json:"receiver_id"`
    Content    string `gorm:"type:text;not null" json:"content"`
    Type       string `gorm:"default:'text'" json:"type"` // "text", "image", "pdf"
    IsRead     bool   `gorm:"default:false" json:"is_read"`
}