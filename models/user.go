package models

import (
	"gorm.io/gorm"
)

// User represents a user in the system.
type User struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey" json:"id"` // Primary key field (auto-incremented ID)
	Username    string `gorm:"type:varchar(100);unique;not null" json:"username"`
	Email       string `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password    string `gorm:"type:varchar(255);not null" json:"password"`
	PhoneNumber string `gorm:"type:varchar(20);not null" json:"phone_number"`
	Token       string `gorm:"type:varchar(255);default:null" json:"token"`

	// Relationships
	Files        []File        `gorm:"foreignKey:UserID" json:"files"`        // Files uploaded by the user
	SharedFiles  []SharedFile  `gorm:"foreignKey:SharedWith" json:"shared_files"` // Files shared with the user
}