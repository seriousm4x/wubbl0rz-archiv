package models

type Game struct {
	UUID      string `gorm:"primaryKey;uniqueIndex;not null" json:"uuid" form:"uuid"`
	Name      string `json:"name" form:"name"`
	BoxartURL string `json:"box_art_url" form:"box_art_url"`
	Clips     []Clip `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"clips" form:"clips"`
}
