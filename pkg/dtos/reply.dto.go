package dtos

import (
	"time"

	"github.com/HudYuSa/comments/database/models"
	"github.com/google/uuid"
)

// because this is a data transfer object
// so everything can be null/empty
// so make everything a pointer except for a type than can be detected by json omitempty as empty value

type ReplyResponse struct {
	ID        *uuid.UUID     `json:"id,omitempty"`
	UserID    *uuid.UUID     `json:"user_id,omitempty"`
	CommentID *uuid.UUID     `json:"comment_id,omitempty"`
	MentionID *uuid.UUID     `json:"mention_id,omitempty"`
	Content   string         `json:"content,omitempty"`
	CreatedAt *time.Time     `json:"created_at,omitempty"`
	User      *UserResponse  `json:"user,omitempty"`
	Post      *PostResponse  `json:"post,omitempty"`
	Mention   *ReplyResponse `json:"mention,omitempty"`
}

type CreateReplyInput struct {
	CommentID uuid.UUID  `json:"comment_id,omitempty"`
	MentionID *uuid.UUID `json:"mention_id,omitempty"`
	Content   string     `json:"content,omitempty" binding:"required"`
}

func ReplyToReplyResponse(reply *models.Reply) *ReplyResponse {
	if reply == nil {
		return nil
	}
	commentResponse := ReplyResponse{
		ID:        CheckNil(reply.ID),
		UserID:    CheckNil(reply.UserID),
		CommentID: CheckNil(reply.CommentID),
		MentionID: reply.MentionID,
		Content:   reply.Content,
		CreatedAt: CheckNil(reply.CreatedAt),
		User:      UserToUserResponse(CheckNil(reply.User)),
		Mention:   ReplyToReplyResponse(reply.Mention),
	}

	return &commentResponse
}
