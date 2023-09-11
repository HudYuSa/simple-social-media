package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	// to use the default uuid generate v4 extension first u need to create/enable the extension from postgresql
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Username       string    `gorm:"not null"`
	Email          string    `gorm:"uniqueIndex;not null"`
	Password       string    `gorm:"not null"`
	ProfilePicture string    `gorm:"not null"`
	CreatedAt      time.Time `gorm:"not null"`
	UpdatedAt      time.Time `gorm:"not null"`
	// VerificationCode   string
	// Verified           bool `gorm:"not null"`
	// PasswordResetToken string
	// PasswordResetAt    time.Time
	// Role               string    `gorm:"type:varchar(255);not null"`
}
