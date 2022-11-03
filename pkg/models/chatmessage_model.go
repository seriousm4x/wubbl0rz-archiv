package models

import (
	"time"
)

type ChatMessage struct {
	ID              string            `gorm:"primarykey"`
	CreatedAt       time.Time         `json:"created_at"`
	UserID          string            `json:"user_id"`
	UserDisplayName string            `json:"user_display_name"`
	UserName        string            `json:"user_name"`
	Message         string            `json:"message"`
	Tags            map[string]string `gorm:"serializer:json" json:"tags"`
}
