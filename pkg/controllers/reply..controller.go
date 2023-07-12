package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/HudYuSa/comments/database/models"
	"github.com/HudYuSa/comments/pkg/dtos"
	"github.com/HudYuSa/comments/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReplyController interface {
	GetReplies(ctx *gin.Context)
	CreateReply(ctx *gin.Context)
	DeleteReply(ctx *gin.Context)
}

type replyController struct {
	DB *gorm.DB
}

func NewReplyController(db *gorm.DB) ReplyController {
	return &replyController{
		DB: db,
	}
}

func (rc *replyController) GetReplies(ctx *gin.Context) {
	commentID := ctx.Param("comment_id")
	replies := []models.Reply{}

	var pageOffset int
	var offsetErr error
	var pageQuery string = ctx.Query("page")

	switch {
	case pageQuery == "", pageQuery < "1":
		pageOffset = 1
	default:
		pageOffset, offsetErr = strconv.Atoi(pageQuery)
	}
	pageOffset = (pageOffset - 1) * 5
	if offsetErr != nil {
		dtos.RespondWithError(ctx, http.StatusBadRequest, offsetErr.Error())
		return
	}

	fmt.Println(pageOffset)
	// result := rc.DB.Preload("User", utils.SelectColumnDB("ID", "Name")).Preload("Replies", func(db *gorm.DB) *gorm.DB {
	// 	return db.Offset((pageOffset - 1) * 10).Limit(10).Order("created_at ASC")
	// }).Preload("Replies.User", utils.SelectColumnDB("ID", "Name")).Preload("Replies.Mention", utils.SelectColumnDB("ID", "UserID")).Preload("Replies.Mention.User", utils.SelectColumnDB("ID", "Name")).Where("id = ?", commentID).First(&comment)
	result := rc.DB.Preload("User", utils.SelectColumnDB("ID", "Name")).
		Preload("Mention", utils.SelectColumnDB("ID", "Name")).
		Where("comment_id = ?", commentID).
		Offset(pageOffset).
		Limit(5).
		Order("created_at ASC").
		Find(&replies)

	if result.Error != nil {
		switch result.Error.Error() {
		default:
			dtos.RespondWithError(ctx, http.StatusBadGateway, result.Error.Error())
		}
		return
	}

	repliesResponse := []dtos.ReplyResponse{}
	for _, reply := range replies {
		repliesResponse = append(repliesResponse, *dtos.ReplyToReplyResponse(&reply))
	}

	dtos.RespondWithJson(ctx, http.StatusOK, repliesResponse)
}

func (rc *replyController) CreateReply(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	payload := dtos.CreateReplyInput{}

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		dtos.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	now := time.Now()
	newReply := models.Reply{
		UserID:    currentUser.ID,
		CommentID: payload.CommentID,
		MentionID: payload.MentionID,
		Content:   payload.Content,
		CreatedAt: now,
	}

	// save the data to database using gorm
	result := rc.DB.Create(&newReply)

	// check for possible error
	if result.Error != nil {
		dtos.RespondWithError(ctx, http.StatusBadGateway, result.Error.Error())
		return
	}

	dtos.RespondWithJson(ctx, http.StatusCreated, dtos.ReplyToReplyResponse(&newReply))
}

func (rc *replyController) DeleteReply(ctx *gin.Context) {
	replyID := ctx.Param("reply_id")

	result := rc.DB.Where("id = ?", replyID).Delete(models.Reply{})
	if result.Error != nil {
		switch result.Error.Error() {
		case "record not found":
			dtos.RespondWithError(ctx, http.StatusNotFound, "cannot delete reply")
		default:
			dtos.RespondWithError(ctx, http.StatusBadGateway, result.Error.Error())
		}
		return
	}

	dtos.RespondWithJson(ctx, http.StatusNoContent, nil)
}
