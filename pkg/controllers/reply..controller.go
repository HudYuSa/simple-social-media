package controllers

import (
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
	comment := models.Comment{}

	var pageOffset int
	var offsetErr error

	if ctx.Query("page") == "" || ctx.Query("page") < "1" {
		pageOffset = 1
	} else {
		pageOffset, offsetErr = strconv.Atoi(ctx.Query("page"))
	}
	if offsetErr != nil {
		dtos.RespondWithError(ctx, http.StatusBadRequest, offsetErr.Error())
		return
	}

	result := rc.DB.Preload("User", utils.SelectColumnDB("ID", "Name")).Preload("Replies").Preload("Replies.User", utils.SelectColumnDB("ID", "Name")).Where("id = ?", commentID).Offset((pageOffset - 1) * 5).Limit(5).Order("created_at ASC").First(&comment)
	// result := rc.DB.Preload("User", utils.SelectColumnDB("ID", "Name")).Preload("Replies.User").Preload("Replies.Mention").Where("comment_id = ?", commentID).Offset((pageOffset - 1) * 5).Limit(5).Order("created_at ASC").First(&comment)
	if result.Error != nil {
		switch result.Error.Error() {
		default:
			dtos.RespondWithError(ctx, http.StatusBadGateway, result.Error.Error())
		}
		return
	}

	// trick untuk infinite scroll
	// buat time stamp pas user nyari, trus cari datanya bedasarkan data sebelum user itu nyari
	// jadi semua data yang dateng saat user sedang nyari bakalan g keliatan
	// usernya harus nyari ulang kalo mau liat data yang baru itu

	dtos.RespondWithJson(ctx, http.StatusOK, dtos.CommentToCommentResponse(&comment))
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
		dtos.RespondWithError(ctx, http.StatusBadGateway, "something bad just happen")
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
