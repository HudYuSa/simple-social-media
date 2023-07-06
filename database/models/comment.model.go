package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	// to use the default uuid generate v4 extension first u need to create/enable the extension from postgresql
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID    uuid.UUID `gorm:"not null"`
	PostID    uuid.UUID `gorm:"not null"`
	Comment   string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	Post      *Post     `gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"Post,omitempty"`
}
