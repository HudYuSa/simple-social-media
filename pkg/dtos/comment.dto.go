package dtos

import (
	"time"

	"github.com/HudYuSa/comments/database/models"
	"github.com/google/uuid"
)

type CommentResponse struct {
	ID              *uuid.UUID        `json:"id,omitempty"`
	UserID          *uuid.UUID        `json:"user_id,omitempty"`
	ParentCommentID *uuid.UUID        `json:"parent_comment_id"`
	PostID          *uuid.UUID        `json:"post_id,omitempty"`
	Content         *string           `json:"content,omitempty"`
	CreatedAt       *time.Time        `json:"created_at,omitempty"`
	ReplyCount      int64             `json:"reply_count"`
	Comments        []CommentResponse `json:"comments,omitempty"`
	Post            *PostResponse     `json:"post,omitempty"`
	User            *UserResponse     `json:"user,omitempty"`
}

type CreateCommentInput struct {
	ParentCommentID *uuid.UUID `json:"parent_comment_id,omitempty"`
	PostID          *uuid.UUID `json:"post_id,omitempty" binding:"required"`
	Content         string     `json:"content,omitempty" binding:"required"`
}

func CommentToCommentResponse(comment *models.Comment) *CommentResponse {
	if comment == nil {
		return nil
	}
	commentsResponse := []CommentResponse{}
	for _, c := range comment.Comments {
		commentsResponse = append(commentsResponse, *CommentToCommentResponse(&c))
	}
	commentResponse := CommentResponse{
		ID:              comment.ID,
		UserID:          comment.UserID,
		PostID:          comment.PostID,
		ParentCommentID: comment.ParentCommentID,
		Content:         comment.Content,
		CreatedAt:       comment.CreatedAt,
		Comments:        commentsResponse,
		Post:            PostToPostResponse(comment.Post),
		User:            UserToUserResponse(comment.User),
	}
	// if comment.ParentCommentID != nil {
	// 	commentResponse.ParentCommentID = comment.ParentCommentID
	// }
	return &commentResponse
}
