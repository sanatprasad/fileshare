package controllers

import (
    "encoding/json"
    "log"
    "net/http"
    "backend/database"
    "backend/models"
    "github.com/go-redis/redis/v8"
    "context"
    "strconv"
)

var ctx = context.Background()

var redisClient *redis.Client

func SetRedisClient(client *redis.Client) {
    redisClient = client
}

func FileMetadataHandler(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(uint)

    cacheKey := "file_metadata:" + strconv.Itoa(int(userID))
    cachedData, err := redisClient.Get(ctx, cacheKey).Result()
    if err == redis.Nil {
        log.Println("Cache miss, fetching from DB")

        var files []models.File
        if err := database.DB.Where("user_id = ?", userID).Find(&files).Error; err != nil {
            log.Printf("Error fetching files: %v\n", err)
            http.Error(w, "Error retrieving files", http.StatusInternalServerError)
            return
        }

        fileData, _ := json.Marshal(files)
        redisClient.Set(ctx, cacheKey, fileData, 0)

        w.Header().Set("Content-Type", "application/json")
        w.Write(fileData)
    } else if err != nil {
        log.Printf("Error accessing cache: %v\n", err)
        http.Error(w, "Error retrieving files", http.StatusInternalServerError)
    } else {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(cachedData))
    }
}
