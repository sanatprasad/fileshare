package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

// File represents the metadata for a file uploaded by a user.
type File struct {
    ID         uuid.UUID   `gorm:"type:uuid;primaryKey;" json:"id"`
    UserID     uuid.UUID   `gorm:"type:uuid;not null;index" json:"user_id"`         // Foreign key to the User
    FileName   string      `gorm:"type:varchar(255);not null" json:"file_name"`     // Name of the file
    Size       int64       `gorm:"not null" json:"size"`                            // File size in bytes
    FileType   string      `gorm:"type:varchar(50)" json:"file_type"`               // Type of the file (e.g., image/jpeg, application/pdf)
    S3URL      string      `gorm:"type:varchar(500);not null" json:"s3_url"`        // URL to access the file in S3 or local storage
    UploadDate time.Time   `gorm:"autoCreateTime" json:"upload_date"`               // Date and time the file was uploaded
    IsPublic   bool        `gorm:"default:false" json:"is_public"`                  // Flag indicating if the file is publicly accessible
    ExpireAt   *time.Time  `json:"expire_at,omitempty"`                            // Optional expiration time for file access (for shared files)

    // Relationships
    User        User         `gorm:"foreignKey:UserID;references:ID" json:"user"` // Reference to the User (many-to-one)
    SharedFiles []SharedFile `gorm:"foreignKey:FileID" json:"shared_files,omitempty"` // Shared files (one-to-many)
}

// BeforeCreate is a GORM hook that sets the UUID before creating a file.
func (f *File) BeforeCreate(tx *gorm.DB) (err error) {
    f.ID = uuid.New()
    return
}