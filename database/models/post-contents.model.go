package models

import "github.com/google/uuid"

type PostContents struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	PostID uuid.UUID `gorm:"not null"`
	Source string    `gorm:"not null"`
	Post   Post      `gorm:"foreignKey:PostID"`
}
