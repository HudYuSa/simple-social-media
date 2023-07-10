package dtos

import (
	"time"

	"github.com/HudYuSa/comments/database/models"
	"github.com/google/uuid"
)

// because this is a data transfer object
// so everything can be null/empty
// so make everything a pointer except for a type than can be detected by json omitempty as empty value

type PostResponse struct {
	ID        *uuid.UUID    `json:"id,omitempty" `
	UserID    *uuid.UUID    `json:"user_id,omitempty" `
	Title     string        `json:"title,omitempty" `
	Photo     string        `json:"photo,omitempty" `
	Content   string        `json:"content,omitempty" `
	CreatedAt *time.Time    `json:"created_at,omitempty"`
	UpdatedAt *time.Time    `json:"updated_at,omitempty"`
	User      *UserResponse `json:"user,omitempty"`
}

type CreatePostInput struct {
	Title   string `json:"title" binding:"required"`
	Photo   string `json:"photo" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type UpdatePostInput struct {
	Title   string `json:"title"`
	Photo   string `json:"photo"`
	Content string `json:"content"`
}

func PostToPostResponse(post *models.Post) *PostResponse {
	if post == nil {
		return nil
	}
	postResponse := PostResponse{
		ID:        CheckNil(post.ID),
		UserID:    CheckNil(post.UserID),
		Title:     post.Title,
		Photo:     post.Photo,
		Content:   post.Content,
		CreatedAt: &post.CreatedAt,
		UpdatedAt: &post.UpdatedAt,
		User:      UserToUserResponse(CheckNil(post.User)),
	}

	return &postResponse
}
