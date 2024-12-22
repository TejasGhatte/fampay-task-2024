package models

import (
	"time"

	"github.com/google/uuid"
)

type Video struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Creator 	  string 		 `gorm:"type:varchar(255)" json:"creator"`
	Title         string         `gorm:"type:varchar(255)" json:"title"`
	Description   string         `gorm:"type:text" json:"description"`
	ThumbnailURLs string`gorm:"type:text[]" json:"thumbnailURLs"`
	PublishedAt     time.Time      `gorm:"default:current_timestamp" json:"createdAt"`
}
