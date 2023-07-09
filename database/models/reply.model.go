package models

import (
	"time"

	"github.com/google/uuid"
)

type Reply struct {
	// to use the default uuid generate v4 extension first u need to create/enable the extension from postgresql
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID    uuid.UUID `gorm:"not null"`
	CommentID uuid.UUID `gorm:"not null"`
	MentionID *uuid.UUID
	Content   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID"`
	Mention   *Reply    `gorm:"foreignKey:MentionID"`
}
