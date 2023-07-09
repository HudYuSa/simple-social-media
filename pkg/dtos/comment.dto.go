package dtos

import (
	"fmt"
	"time"

	"github.com/HudYuSa/comments/database/models"
	"github.com/google/uuid"
)

type CommentResponse struct {
	ID         uuid.UUID        `json:"id,omitempty"`
	UserID     uuid.UUID        `json:"user_id,omitempty"`
	PostID     uuid.UUID        `json:"post_id,omitempty"`
	Content    string           `json:"content,omitempty"`
	CreatedAt  time.Time        `json:"created_at,omitempty"`
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
		ID:        comment.ID,
		UserID:    comment.UserID,
		PostID:    comment.PostID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
		User:      UserToUserResponse(CheckNil(comment.User)),
		Post:      PostToPostResponse(CheckNil(comment.Post)),
		Replies:   repliesResponse,
	}
	// if comment.ParentCommentID != nil {
	// 	commentResponse.ParentCommentID = comment.ParentCommentID
	// }
	return &commentResponse
}
