package controllers

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"
    "backend/database"
    "backend/models"
    "github.com/gorilla/mux"
    "github.com/google/uuid"
)

// ShareFileHandler generates a public shareable link for a file
func ShareFileHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    fileID := vars["file_id"]

    userID := r.Context().Value("userID").(uint)

    // Fetch the file from the database
    var file models.File
    if err := database.DB.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
        log.Printf("Error fetching file: %v\n", err)
        http.Error(w, "File not found", http.StatusNotFound)
        return
    }

    // Create a new shared file entry with optional expiration
    sharedFile := models.SharedFile{
        ID:        uuid.New(),
        FileID:    file.ID,
        UserID:    userID,
        ExpiresAt: time.Now().Add(24 * time.Hour), // Example: 24-hour expiration
    }

    if err := database.DB.Create(&sharedFile).Error; err != nil {
        log.Printf("Error creating shared file: %v\n", err)
        http.Error(w, "Error sharing file", http.StatusInternalServerError)
        return
    }

    // Generate the public URL
    publicURL := fmt.Sprintf("https://yourdomain.com/share/%s", sharedFile.ID.String())

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "public_url": publicURL,
    })
}
