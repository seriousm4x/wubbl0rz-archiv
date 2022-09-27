package models

type Emote struct {
	ID       string `gorm:"primaryKey;uniqueIndex;not null" json:"id" form:"id"`
	Name     string `json:"name" form:"name"`
	URL      string `json:"url" form:"url"`
	Provider string `json:"provider" form:"provider"`
	Outdated bool   `json:"-" form:"-"`
}
