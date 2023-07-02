package dtos

import (
	"time"

	"github.com/HudYuSa/mod-name/database/models"
	"github.com/google/uuid"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id,omitempty" form:"id"`
	Name      string    `json:"name,omitempty" form:"name"`
	Email     string    `json:"email,omitempty" form:"email"`
	Role      string    `json:"role,omitempty" form:"role"`
	Photo     string    `json:"photo,omitempty" form:"photo"`
	Provider  string    `json:"provider" form:"provider"`
	CreatedAt time.Time `json:"created_at" form:"created_at"`
	UpdatedAt time.Time `json:"updated_at" form:"updated_at"`
}

type SignUpInput struct {
	Name            string `json:"name" binding:"required" form:"name"`
	Email           string `json:"email" binding:"required,email" form:"email"`
	Password        string `json:"password" binding:"required,min=8" form:"password"`
	PasswordConfirm string `json:"password_confirm" binding:"required" form:"password_confirm"`
	Photo           string `json:"photo" form:"photo"`
}

type SignInInput struct {
	Email    string `json:"email" binding:"required,email" form:"email"`
	Password string `json:"password" binding:"required" form:"password"`
}

// ? ForgotPassowrdInput struct
type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required,email"`
}

// ? ResetPasswordInput struct
type ResetPasswordInput struct {
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"password_confirm" biding:"required"`
}

func DBUserToUserResponse(dbUser *models.User) UserResponse {
	return UserResponse{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		Photo:     dbUser.Photo,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}
}
