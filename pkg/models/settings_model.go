package models

import (
	"time"
)

type Settings struct {
	ID                  int `gorm:"primaryKey;uniqueIndex;not null"`
	BroadcasterId       string
	TtvClientId         string
	TtvClientSecret     string
	TtvBearerToken      string
	TtvBearerExpireDate time.Time
	DateVodsUpdate      time.Time
	DateEmotesUpdate    time.Time
	IsLive              bool
}
