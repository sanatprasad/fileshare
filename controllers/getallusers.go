package controllers

import (
    "encoding/json"
    "log"
    "net/http"
    "backend/database"
    "backend/models" 
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
    db := database.GetDB() 
    var users []models.User

    err := db.Find(&users).Error
    if err != nil {
        log.Println("Error fetching users:", err)
        http.Error(w, "Error fetching users", http.StatusInternalServerError)
        return
    }

    response := map[string]interface{}{
        "status": "success",
        "data":   users,
    }

    w.Header().Set("Content-Type", "application/json")

    if err := json.NewEncoder(w).Encode(response); err != nil {
        log.Println("Error encoding response:", err)
        http.Error(w, "Error encoding response", http.StatusInternalServerError)
    }
}
