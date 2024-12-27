package utils

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Paginator(c *fiber.Ctx) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		pageStr := c.Query("page", "1")
		limitStr := c.Query("limit", "10")

		page, err := strconv.Atoi(pageStr)
		if err != nil {
			log.Println("Failed to Paginate due to integer conversion.")
			return db
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			log.Println("Failed to Paginate due to integer conversion.")
			return db
		}

		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}

func CursorPaginateVideos(c *fiber.Ctx) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		cursorStr := c.Query("cursor", "")
		limitStr := c.Query("limit", "10")

		limit, err := strconv.Atoi(limitStr)
        if err != nil || limit <= 0 {
            limit = 10
        }

        db = db.Limit(limit)
		cursorStr = strings.ReplaceAll(cursorStr, " ", "+")
		cursor, err := time.Parse(time.RFC3339, cursorStr)
		if err != nil {
			fmt.Printf("Invalid Cursor: %s", err.Error())
			return db
		}

		if cursorStr != "" {
			db = db.Where("published_at < ?", cursor)
		}

		return db
	}
}