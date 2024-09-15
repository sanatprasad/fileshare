package controllers

import (
    "bytes"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/google/uuid"
    "gorm.io/gorm"

    "backend/models"
)

var (
    // AWS configurations are now loaded from environment variables
    awsAccessKeyID     = os.Getenv("AWS_ACCESS_KEY_ID")
    awsSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
    region             = os.Getenv("AWS_REGION")
    bucket             = os.Getenv("S3_BUCKET_NAME")

    db *gorm.DB
)

// SetDB initializes the database connection
func SetDB(database *gorm.DB) {
    db = database
}

// UploadFile handles file uploads and saves metadata
func UploadFile(w http.ResponseWriter, r *http.Request) {
    // Retrieve the user ID from the JWT middleware
    userID := r.Header.Get("UserID")
    if userID == "" {
        http.Error(w, "Unauthorized: No user ID found", http.StatusUnauthorized)
        return
    }

    // Parse form with multipart files
    err := r.ParseMultipartForm(10 << 20) // Max upload size: 10MB
    if err != nil {
        http.Error(w, "Error parsing form", http.StatusBadRequest)
        return
    }

    file, handler, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "Error retrieving the file", http.StatusBadRequest)
        return
    }
    defer file.Close()

    // Create an AWS session using environment variables
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String(region),
        Credentials: credentials.NewStaticCredentials(
            awsAccessKeyID,
            awsSecretAccessKey,
            "",
        ),
    })
    if err != nil {
        http.Error(w, "Error creating AWS session", http.StatusInternalServerError)
        log.Println("AWS Session Error:", err)
        return
    }

    svc := s3.New(sess)

    // Read the file content into a buffer
    var buf bytes.Buffer
    _, err = io.Copy(&buf, file)
    if err != nil {
        http.Error(w, "Failed to read file", http.StatusInternalServerError)
        return
    }

    // Generate a unique key for the file in S3 using user ID and UUID
    uniqueKey := fmt.Sprintf("%s/%s-%s", userID, uuid.New().String(), handler.Filename)

    // Upload file to S3
    _, err = svc.PutObject(&s3.PutObjectInput{
        Bucket:      aws.String(bucket),
        Key:         aws.String(uniqueKey),
        Body:        bytes.NewReader(buf.Bytes()),
        ContentType: aws.String(handler.Header.Get("Content-Type")),
    })
    if err != nil {
        log.Println("AWS S3 Error:", err) // Log the detailed error
        http.Error(w, "Failed to upload file to S3", http.StatusInternalServerError)
        return
    }

    // Save file metadata to PostgreSQL
    fileSize := handler.Size
    fileType := handler.Header.Get("Content-Type")
    s3URL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, uniqueKey)

    // Create a new File record
    fileRecord := models.File{
        ID:         uuid.New(),
        UserID:     uuid.MustParse(userID), // Use the actual user ID from JWT
        FileName:   handler.Filename,
        Size:       fileSize,
        FileType:   fileType,
        S3URL:      s3URL,
        UploadDate: time.Now(),
        IsPublic:   false, // Set this according to your requirements
    }

    // Save metadata to the database
    if err := db.Create(&fileRecord).Error; err != nil {
        http.Error(w, "Failed to save file metadata", http.StatusInternalServerError)
        return
    }

    // Return success message with the file URL
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "File uploaded successfully: %s\nFile URL: %s\n", handler.Filename, s3URL)
}
