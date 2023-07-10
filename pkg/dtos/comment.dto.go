package dtos

import (
	"fmt"
	"time"

	"github.com/HudYuSa/comments/database/models"
	"github.com/google/uuid"
)

// because this is a data transfer object
// so everything can be null/empty
// so make everything a pointer except for a type than can be detected by json omitempty as empty value

type CommentResponse struct {
	ID         *uuid.UUID       `json:"id,omitempty"`
	UserID     *uuid.UUID       `json:"user_id,omitempty"`
	PostID     *uuid.UUID       `json:"post_id,omitempty"`
	Content    string           `json:"content,omitempty"`
	CreatedAt  *time.Time       `json:"created_at,omitempty"`
	ReplyCount int64            `json:"reply_count"`
	User       *UserResponse    `json:"user,omitempty"`
	Post       *PostResponse    `json:"post,omitempty"`
	Replies    []ReplyResponse  `json:"replies,omitempty"`
	Mention    *CommentResponse `json:"mention,omitempty"`
}

type CreateCommentInput struct {
	PostID  *uuid.UUID `json:"post_id,omitempty" binding:"required"`
	Content string     `json:"content,omitempty" binding:"required"`
}

func CommentToCommentResponse(comment *models.Comment) *CommentResponse {
	if comment == nil {
		return nil
	}
	var repliesResponse []ReplyResponse

	for _, r := range comment.Replies {
		replyResponse := *ReplyToReplyResponse(&r)
		repliesResponse = append(repliesResponse, replyResponse)
	}
	fmt.Println(repliesResponse)

	commentResponse := CommentResponse{
		ID:        CheckNil(comment.ID),
		UserID:    CheckNil(comment.UserID),
		PostID:    CheckNil(comment.PostID),
		Content:   comment.Content,
		CreatedAt: CheckNil(comment.CreatedAt),
		User:      UserToUserResponse(CheckNil(comment.User)),
		Post:      PostToPostResponse(CheckNil(comment.Post)),
		Replies:   repliesResponse,
	}

	return &commentResponse
}
