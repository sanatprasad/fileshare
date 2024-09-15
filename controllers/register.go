package controllers

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/dgrijalva/jwt-go"
    "github.com/joho/godotenv"
    "backend/database"
    "backend/models"
    "golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
    Username        string `json:"username"`
    Email           string `json:"email"`
    Password        string `json:"password"`
    ConfirmPassword string `json:"confirm_password"`
    PhoneNumber     string `json:"phone_number"`
}

// Load environment variables from .env file
func init() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }
}

// Retrieve JWT secret key from environment variable
var jwtSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("RegisterHandler called")

    if r.Method == http.MethodPost {
        var req RegisterRequest
        log.Println("Decoding JSON request body")

        err := json.NewDecoder(r.Body).Decode(&req)
        if err != nil {
            log.Printf("Error decoding JSON: %v\n", err)
            http.Error(w, "Unable to parse JSON", http.StatusBadRequest)
            return
        }

        if req.Password != req.ConfirmPassword {
            log.Println("Passwords do not match")
            http.Error(w, "Passwords do not match", http.StatusBadRequest)
            return
        }

        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
        if err != nil {
            log.Printf("Error hashing password: %v\n", err)
            http.Error(w, "Error creating user", http.StatusInternalServerError)
            return
        }

        user := models.User{
            Username:    req.Username,
            Email:       req.Email,
            Password:    string(hashedPassword),
            PhoneNumber: req.PhoneNumber,
        }

        result := database.DB.Create(&user)
        if result.Error != nil {
            log.Printf("Error creating user: %v\n", result.Error)
            http.Error(w, "Error creating user", http.StatusInternalServerError)
            return
        }

        token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
            "sub":   user.ID, 
            "exp":   time.Now().Add(time.Hour * 24).Unix(),
            "email": user.Email,
        })

        tokenString, err := token.SignedString(jwtSecretKey)
        if err != nil {
            log.Printf("Error generating token: %v\n", err)
            http.Error(w, "Error creating token", http.StatusInternalServerError)
            return
        }

        user.Token = tokenString
        if err := database.DB.Save(&user).Error; err != nil {
            log.Printf("Error saving user token: %v\n", err)
            http.Error(w, "Error updating user token", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        response := map[string]string{
            "message": "Registration successful",
            "token":   tokenString,
        }
        json.NewEncoder(w).Encode(response)
    } else {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
    }
}
