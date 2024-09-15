package controllers

import (
    "encoding/json"
    "net/http"
    "backend/database"
    "backend/models"
    "github.com/dgrijalva/jwt-go"
    "golang.org/x/crypto/bcrypt"
    "time"
)

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        var req LoginRequest
        err := json.NewDecoder(r.Body).Decode(&req)
        if err != nil {
            http.Error(w, "Unable to parse JSON", http.StatusBadRequest)
            return
        }

        var user models.User
        result := database.DB.Where("username = ?", req.Username).First(&user)
        if result.Error != nil {
            http.Error(w, "Invalid username or password", http.StatusUnauthorized)
            return
        }

        // Compare the hashed password with the one provided in the request
        err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
        if err != nil {
            http.Error(w, "Invalid username or password", http.StatusUnauthorized)
            return
        }

        // Generate JWT
        token, err := generateJWT(user)
        if err != nil {
            http.Error(w, "Error generating JWT", http.StatusInternalServerError)
            return
        }

        // Return the JWT to the user
        w.Header().Set("Content-Type", "application/json")
        response := map[string]string{"message": "Login successful", "token": token}
        json.NewEncoder(w).Encode(response)
    } 
	else {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
    }
}

func generateJWT(user models.User) (string, error) { 
    claims := jwt.MapClaims{
        "sub": user.ID,
        "exp": time.Now().Add(time.Hour * 24).Unix(), // Token expiration time
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte("secret")) // Use a secure secret key
}
