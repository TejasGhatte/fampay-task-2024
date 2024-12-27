package controllers

import (
	"time"

	"github.com/TejasGhatte/fampay-task-2024/initializers"
	"github.com/TejasGhatte/fampay-task-2024/models"
	"github.com/TejasGhatte/fampay-task-2024/utils"
	"github.com/gofiber/fiber/v2"
)

func GetVideos(c *fiber.Ctx) error {
	paginatedDB := utils.CursorPaginateVideos(c)(initializers.DB)

	var videos []models.Video
	if err := paginatedDB.Order("published_at DESC").Find(&videos).Error; err != nil {
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: "Failed to get videos"}
	}

	var nextCursor string
	if len(videos) > 0 {
		nextCursor = videos[len(videos)-1].PublishedAt.Format(time.RFC3339)
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"message": "Videos fetched successfully",
		"videos":   videos,
		"next_cursor": nextCursor,
	})
}