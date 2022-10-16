package models

import "time"

type Clip struct {
	UUID        string     `gorm:"primaryKey;uniqueIndex" json:"uuid" form:"uuid"`
	Title       string     `json:"title" form:"title"`
	Duration    int        `json:"duration" form:"duration"`
	Date        *time.Time `json:"date" form:"date" time_format:"2006-01-02T15:04:05.000Z"`
	Filename    string     `json:"filename" form:"filename"`
	Resolution  string     `json:"resolution" form:"resolution"`
	Fps         float32    `json:"fps" form:"fps"`
	Size        int        `json:"size" form:"size"`
	Viewcount   int        `gorm:"default:0" json:"viewcount" form:"viewcount"`
	VodOffset   int        `json:"vod_offset" form:"vod_offset"`
	CreatorUUID string     `json:"creator_uuid" form:"creator_uuid"`
	Creator     Creator    `gorm:"foreignKey:CreatorUUID;references:UUID" form:"creator" json:"creator"`
	GameUUID    string     `json:"game_uuid" form:"game_uuid"`
	Game        Game       `gorm:"foreignKey:GameUUID;references:UUID" json:"game" form:"game"`
	VodUUID     string     `json:"vod_uuid" form:"vod_uuid"`
	Vod         Vod        `gorm:"foreignKey:VodUUID;references:UUID" json:"vod" form:"vod"`
}
