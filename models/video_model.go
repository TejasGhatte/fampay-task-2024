package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Video struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Creator 	  string 		 `gorm:"type:varchar(255)" json:"creator"`
	Title         string         `gorm:"type:varchar(255)" json:"title"`
	Description   string         `gorm:"type:text" json:"description"`
	ThumbnailURLs pq.StringArray `gorm:"type:text[]" json:"thumbnailURLs"`
	CreatedAt     time.Time      `gorm:"default:current_timestamp" json:"createdAt"`
}
