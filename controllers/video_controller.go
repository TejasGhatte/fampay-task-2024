package controllers

import (
	"log"
	"strconv"

	"github.com/TejasGhatte/fampay-task-2024/initializers"
	"github.com/TejasGhatte/fampay-task-2024/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetVideos(c *fiber.Ctx) error {
	paginatedDB := Paginator(c)(initializers.DB)

	var videos []models.Video
	if err := paginatedDB.Order("published_at DESC").Find(&videos).Error; err != nil {
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: "Failed to get videos"}
	}
	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"message": "Videos fetched successfully",
		"data":   videos,
	})
}

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