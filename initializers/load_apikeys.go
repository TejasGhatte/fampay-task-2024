package initializers

import (
	"encoding/json"
	"log"
	"os"
)

type APIKeys struct {
    YouTubeAPIKeys []string `json:"youtube_api_keys"`
}

var ApiKeys []string

func LoadAPIKeys() []string {
    data, err := os.ReadFile("config/apikeys.json")
    if err != nil {
        log.Fatalf("Error reading config file: %v", err)
    }

    var config APIKeys
    if err := json.Unmarshal(data, &config); err != nil {
        log.Fatalf("Error parsing config file: %v", err)
    }

    ApiKeys = config.YouTubeAPIKeys
	if len(ApiKeys) == 0 {
		log.Fatal("No YouTube API keys found in the configuration file.")
	}
	return ApiKeys
}
