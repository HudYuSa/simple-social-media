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
	dataTimeString := ctx.Query("data-time")

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

	// if the request query don't have datatime then just use current time as datatime
	if dataTimeString == "" {
		dataTimeString = time.Now().Format(time.DateTime)
	}

	dataTime, err := time.Parse(time.DateTime, dataTimeString)
	if err != nil {
		dtos.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// result := cc.DB.Preload("User", utils.SelectColumnDB("ID", "Name")).Preload("Replies").Where("post_id = ? AND created_at < ?", postID, dataTime).Find(&comments).Offset((pageOffset - 1) * 5).Limit(5).Order("created_at DESC")
	// result := cc.DB.Raw("SELECT c.*, COUNT(r.id) as reply_count FROM comments c LEFT JOIN replies r ON c.id = r.comment_id WHERE post_id = ? AND c.created_at < ? GROUP BY c.id ORDER BY created_at DESC OFFSET ? LIMIT ? ", postID, dataTime, pageOffset, 5).Scan(&comments)

	result := cc.DB.
		Preload("User", utils.SelectColumnDB("ID", "Name")).
		Table("comments c").
		Select("c.*, COUNT(r.id) as reply_count").
		Where("post_id = ? AND c.created_at < ? ", postID, dataTime).
		Joins("LEFT JOIN replies r on r.comment_id = c.id").
		Group("c.id").
		Order("c.created_at DESC").
		Offset(pageOffset).
		Limit(5).
		Find(&comments)

	if result.Error != nil {
		switch result.Error.Error() {
		default:
			dtos.RespondWithError(ctx, http.StatusBadGateway, result.Error.Error())
		}
		return
	}

	commentsResponse := []*dtos.CommentResponse{}
	for _, comment := range comments {
		commentsResponse = append(commentsResponse, dtos.CommentToCommentResponse(&comment))
	}

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
		dtos.RespondWithError(ctx, http.StatusBadGateway, result.Error.Error())
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
