package models

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID    *uuid.UUID `gorm:"not null"`
	Title     string     `gorm:"not nul"`
	Photo     string     `gorm:"not null"`
	Content   string     `gorm:"type:text;not null"`
	CreatedAt time.Time  `gorm:"not null"`
	UpdatedAt time.Time  `gorm:"not null"`
	User      User       `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
