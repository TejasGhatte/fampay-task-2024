package routines

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/TejasGhatte/fampay-task-2024/initializers"
	"github.com/TejasGhatte/fampay-task-2024/models"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var (
	query      = flag.String("query", "football", "Search term")
	maxResults = flag.Int64("max-results", 100, "Max YouTube results")
)

const developerKey = "API_KEY"

func FetchVideos(c *fiber.Ctx) error {
	flag.Parse()

	// Initialize YouTube service
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(developerKey))
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	// Fetch videos from YouTube API
	err = fetchAndStoreVideos(service)
	if err != nil {
		log.Printf("Error fetching and storing videos: %v", err)
	}
	fmt.Println("Videos fetched and stored successfully")

	return nil
}

func fetchAndStoreVideos(service *youtube.Service) error {
	// Make the API call to YouTube
	call := service.Search.List([]string{"id", "snippet"}).
		Q(*query).
		Type("video").
		MaxResults(*maxResults).
		Order("date")// Order by date to get the latest videos

	response, err := call.Do()
	if err != nil {
		return fmt.Errorf("error calling Search.List: %v", err)
	}

	fmt.Println(response)

	for _, item := range response.Items {
		publishedAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
		if err != nil {
			return fmt.Errorf("error parsing publishedAt: %v", err)
		}
		video := models.Video{
			Title : item.Snippet.Title,
			Description :item.Snippet.Description,
			PublishedAt : publishedAt,
			ThumbnailURLs : item.Snippet.Thumbnails.Default.Url,
		}
		fmt.Println(video.Title)

		// Insert video details into the database
		if err := initializers.DB.Create(&video).Error; err != nil {
			return fmt.Errorf("error inserting video into database: %v", err)
		}
	}

	return nil
}