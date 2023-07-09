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

type CommentController interface {
	GetComments(ctx *gin.Context)
	CreateComment(ctx *gin.Context)
	DeleteComment(ctx *gin.Context)
}

type commentController struct {
	DB *gorm.DB
}

func NewCommentController(db *gorm.DB) CommentController {
	return &commentController{DB: db}
}

func (cc *commentController) GetComments(ctx *gin.Context) {
	postID := ctx.Param("post_id")
	comments := []models.Comment{}

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

	result := cc.DB.Preload("User", utils.SelectColumnDB("ID", "Name")).Where("post_id = ?", postID).Offset((pageOffset - 1) * 5).Limit(5).Order("created_at DESC").Find(&comments)
	if result.Error != nil {
		switch result.Error.Error() {
		default:
			dtos.RespondWithError(ctx, http.StatusBadGateway, result.Error.Error())
		}
		return
	}

	commentsResponse := []*dtos.CommentResponse{}
	for _, comment := range comments {
		var replyCount int64
		// for each comment find and count how many comment have the same
		// parent_coment_id with the current comment id
		result := cc.DB.Where("comment_id = ?", comment.ID).Find(&models.Reply{}).Count(&replyCount)
		if result.Error != nil {
			switch result.Error.Error() {
			default:
				dtos.RespondWithError(ctx, http.StatusBadGateway, result.Error.Error())
			}
			return
		}
		commentResponse := dtos.CommentToCommentResponse(&comment)
		commentResponse.ReplyCount = replyCount

		commentsResponse = append(commentsResponse, commentResponse)
	}
	// trick untuk infinite scroll
	// buat time stamp pas user nyari, trus cari datanya bedasarkan data sebelum user itu nyari
	// jadi semua data yang dateng saat user sedang nyari bakalan g keliatan
	// usernya harus nyari ulang kalo mau liat data yang baru itu

	dtos.RespondWithJson(ctx, http.StatusOK, commentsResponse)
}

func (cc *commentController) CreateComment(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	payload := dtos.CreateCommentInput{}

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		dtos.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	now := time.Now()
	newComment := models.Comment{
		UserID:    currentUser.ID,
		PostID:    *payload.PostID,
		Content:   payload.Content,
		CreatedAt: now,
	}

	// save the data to database using gorm
	result := cc.DB.Create(&newComment)

	// check for possible error
	if result.Error != nil {
		dtos.RespondWithError(ctx, http.StatusBadGateway, "something bad just happen")
		return
	}

	dtos.RespondWithJson(ctx, http.StatusCreated, dtos.CommentToCommentResponse(&newComment))
}

func (cc *commentController) DeleteComment(ctx *gin.Context) {
	commentID := ctx.Param("comment_id")

	result := cc.DB.Where("id = ?", commentID).Delete(models.Comment{})
	if result.Error != nil {
		switch result.Error.Error() {
		case "record not found":
			dtos.RespondWithError(ctx, http.StatusNotFound, "cannot delete comment")
		default:
			dtos.RespondWithError(ctx, http.StatusBadGateway, result.Error.Error())
		}
		return
	}

	dtos.RespondWithJson(ctx, http.StatusNoContent, nil)
}
