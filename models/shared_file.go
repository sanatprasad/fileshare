package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

// SharedFile represents a shared file with additional details.
type SharedFile struct {
    ID        uuid.UUID `gorm:"type:uuid;primaryKey;" json:"id"`
    FileID    uuid.UUID `gorm:"type:uuid;not null;index" json:"file_id"`    // Foreign key to the File
    SharedWith uuid.UUID `gorm:"type:uuid;not null;index" json:"shared_with"` // Foreign key to the User
    ExpiryDate *time.Time `json:"expiry_date,omitempty"`                    // Optional expiration date for sharing

    // Relationships
    File File `gorm:"foreignKey:FileID;references:ID" json:"file"` // Reference to the File (many-to-one)
    User User `gorm:"foreignKey:SharedWith;references:ID" json:"user"` // Reference to the User (many-to-one)
}

// BeforeCreate is a GORM hook that sets the UUID before creating a shared file.
func (s *SharedFile) BeforeCreate(tx *gorm.DB) (err error) {
    s.ID = uuid.New()
    return
}
