package models

type Creator struct {
	UUID  string `gorm:"primaryKey;uniqueIndex;not null" json:"uuid"`
	Name  string `json:"name"`
	Clips []Clip `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"clips"`
}
