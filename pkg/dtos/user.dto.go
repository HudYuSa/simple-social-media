package dtos

import (
	"time"

	"github.com/HudYuSa/comments/database/models"
	"github.com/google/uuid"
)

// because this is a data transfer object
// so everything can be null/empty
// so make everything a pointer except for a type than can be detected by json omitempty as empty value

type UserResponse struct {
	ID        *uuid.UUID `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	Email     string     `json:"email,omitempty"`
	Role      string     `json:"role,omitempty"`
	Photo     string     `json:"photo,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type SignUpInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type SignInInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func UserToUserResponse(user *models.User) *UserResponse {
	if user == nil {
		return nil
	}
	return &UserResponse{
		ID:        CheckNil(user.ID),
		Name:      user.Name,
		Email:     user.Email,
		Photo:     user.Photo,
		CreatedAt: CheckNil(user.CreatedAt),
		UpdatedAt: CheckNil(user.UpdatedAt),
	}
}
