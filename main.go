package main

import (
    "backend/controllers"
    "context"
    "github.com/go-redis/redis/v8"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
    "backend/routes"
)

func main() {
    // Load environment variables from the system
    redisAddr := os.Getenv("REDIS_HOST") + ":6379"
    redisPassword := os.Getenv("REDIS_PASSWORD") // Set in your .env file if required
    serverAddr := os.Getenv("SERVER_ADDRESS")    // Define SERVER_ADDRESS in your .env file (e.g., ":8081")

    if serverAddr == "" {
        serverAddr = ":8081" // Default to port 8081 if not set
    }

    // Initialize Redis client using environment variables
    redisClient := redis.NewClient(&redis.Options{
        Addr:     redisAddr,
        Password: redisPassword, // No password if empty
        DB:       0,             // Default DB
    })

    // Ping Redis to verify the connection
    _, err := redisClient.Ping(context.Background()).Result()
    if err != nil {
        log.Fatalf("Failed to connect to Redis: %v\n", err)
    }

    // Set Redis client in controllers
    controllers.SetRedisClient(redisClient)

    // Set up the router (mux)
    router := routes.SetupRoutes()

    // Starting the HTTP server
    srv := &http.Server{
        Addr:         serverAddr, // Server address from environment variable
        Handler:      router,     // Pass the router here
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  15 * time.Second,
    }

    // Channel to listen for OS signals for graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

    // Run server in a goroutine to allow graceful shutdown
    go func() {
        log.Printf("Server started on port %s\n", serverAddr)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server failed to start: %v\n", err)
        }
    }()

    // Block until we receive an OS signal
    <-quit
    log.Println("Shutting down server...")

    // Gracefully shutdown the server
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v\n", err)
    }

    log.Println("Server exited gracefully")
}
