package middlewares

import (
    "context"
    "net/http"
    "strings"
    "log"

    "github.com/dgrijalva/jwt-go"
    "backend/models"
    "backend/database"
    "errors"
)

var jwtSecretKey = []byte("secret") // Replace with your actual secret key

// Define the error
var ErrInvalidToken = errors.New("invalid token")

// Context key type for user ID
type contextKey string

const userIDKey contextKey = "userID"

// JWTAuthMiddleware is a middleware for JWT authentication
func JWTAuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tokenString := r.Header.Get("Authorization")
        if tokenString == "" {
            http.Error(w, "Missing token", http.StatusUnauthorized)
            return
        }

        // Remove "Bearer " prefix if it exists
        if strings.HasPrefix(tokenString, "Bearer ") {
            tokenString = strings.TrimPrefix(tokenString, "Bearer ")
        }

        // Parse the JWT token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            // Check the signing method
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, ErrInvalidToken
            }
            return jwtSecretKey, nil
        })
        if err != nil || !token.Valid {
            log.Printf("Invalid token: %v\n", err)
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        // Extract user ID from token claims
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok || !token.Valid {
            log.Printf("Invalid token claims: %v\n", err)
            http.Error(w, "Invalid token claims", http.StatusUnauthorized)
            return
        }

        userID, ok := claims["sub"].(float64) // User ID is expected to be a float64
        if !ok {
            log.Println("User ID missing in token claims")
            http.Error(w, "User ID missing in token claims", http.StatusUnauthorized)
            return
        }

        // Convert userID to uint if necessary
        userIDUint := uint(userID)

        // Optionally fetch user from database to verify existence
        var user models.User
        if err := database.DB.First(&user, "id = ?", userIDUint).Error; err != nil {
            log.Printf("User not found: %v\n", err)
            http.Error(w, "User not found", http.StatusUnauthorized)
            return
        }

        // Set the user ID in the request context
        ctx := context.WithValue(r.Context(), userIDKey, userIDUint)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Retrieve user ID from context
func GetUserIDFromContext(ctx context.Context) (uint, bool) {
    userID, ok := ctx.Value(userIDKey).(uint)
    return userID, ok
}