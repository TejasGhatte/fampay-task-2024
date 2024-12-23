package routines

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/TejasGhatte/fampay-task-2024/initializers"
	"github.com/TejasGhatte/fampay-task-2024/models"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var (
	query      = flag.String("query", "football", "Search term")
	maxResults = flag.Int64("max-results", 25, "Max YouTube results")
)

var apiKeys = initializers.LoadAPIKeys()
var (
    currentKeyIndex = 0
    keyMutex        sync.Mutex
)

func FetchVideos(){
	startTime := time.Now()
	flag.Parse()

	// Initialize YouTube service
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(getCurrentAPIKey()))
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	// Fetch videos from YouTube API
	err = fetchAndStoreVideos(service)
	if err != nil {
		log.Printf("Error fetching and storing videos: %v", err)
	}

	endTime := time.Now()
	fmt.Println("Time taken to fetch and store videos:", endTime.Sub(startTime))
	fmt.Println("Videos fetched and stored successfully")
}

func fetchAndStoreVideos(service *youtube.Service) error {
	lastFetchTime := time.Now().Add(-24 * time.Hour)
	var lastVideo models.Video
	if err := initializers.DB.Order("published_at desc").First(&lastVideo).Error; err == nil {
		lastFetchTime = lastVideo.PublishedAt
	}

	const maxRetries = 5
    retries := 0

	for {
		if retries >= maxRetries {
			return fmt.Errorf("maximum retries exceeded. Exiting")
		}
		// Make the API call to YouTube
		call := service.Search.List([]string{"id", "snippet"}).
		Q(*query).
		Type("video").
		MaxResults(*maxResults).
		Order("date").
		PublishedAfter(lastFetchTime.Format(time.RFC3339))
		// Order by date to get the latest videos

		response, err := call.Do()
		if err != nil {
			if isQuotaExceededError(err) {
				fmt.Println("Quota exceeded. Retrying...")
				if !rotateAPIKey() {
					return fmt.Errorf("all API keys exhausted")
				}
				// Reinitialize the service with the new key
				service, err = youtube.NewService(context.Background(), option.WithAPIKey(getCurrentAPIKey()))
				if err != nil {
					return fmt.Errorf("error reinitializing YouTube client: %v", err)
				}
				retries++
				continue
			}
			return fmt.Errorf("error calling Search.List: %v", err)
		}

		var videos []models.Video
		for _, item := range response.Items {
			// Check for duplicates
			existingVideo := models.Video{}
			if err := initializers.DB.Where("video_id = ?", item.Id.VideoId).First(&existingVideo).Error; err == nil {
				continue
			}
		
			publishedAt, _ := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
			video := models.Video{
				Creator:      item.Snippet.ChannelTitle,
				VideoID:      item.Id.VideoId,
				Title:        item.Snippet.Title,
				Description:  item.Snippet.Description,
				PublishedAt:  publishedAt,
				ThumbnailURLs: []string{item.Snippet.Thumbnails.Default.Url},
			}
		
			videos = append(videos, video)
		}
		
		if len(videos) > 0 {
			if err := initializers.DB.Create(&videos).Error; err != nil {
				return fmt.Errorf("error inserting videos into database: %v", err)
			}
		}
		break
	}

	return nil
}

func getCurrentAPIKey() string {
    keyMutex.Lock()
    defer keyMutex.Unlock()
    return apiKeys[currentKeyIndex]
}

func rotateAPIKey() bool {
    keyMutex.Lock()
    defer keyMutex.Unlock()
    currentKeyIndex++
    return currentKeyIndex < len(apiKeys)
}

func isQuotaExceededError(err error) bool {
    apiErr, ok := err.(*googleapi.Error)
    if !ok {
        return false
    }
    for _, e := range apiErr.Errors {
        if e.Reason == "quotaExceeded" {
            return true
        }
    }
    return false
}