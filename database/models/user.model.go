package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	// to use the default uuid generate v4 extension first u need to create/enable the extension from postgresql
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Password  string    `GORM:"not null"`
	Photo     string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	// VerificationCode   string
	// Verified           bool `gorm:"not null"`
	// PasswordResetToken string
	// PasswordResetAt    time.Time
	// Role               string    `gorm:"type:varchar(255);not null"`
}
