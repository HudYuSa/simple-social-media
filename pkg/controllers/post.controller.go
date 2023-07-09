package controllers

import (
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/HudYuSa/comments/database/models"
	"github.com/HudYuSa/comments/pkg/dtos"
	"github.com/HudYuSa/comments/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostController interface {
	GetPostByID(ctx *gin.Context)
	GetAllPost(ctx *gin.Context)
	CreatePost(ctx *gin.Context)
	UpdatePost(ctx *gin.Context)
	DeletePost(ctx *gin.Context)
}

type postController struct {
	DB *gorm.DB
}

func NewPostController(db *gorm.DB) PostController {
	return &postController{
		DB: db,
	}
}

func (pc *postController) GetPostByID(ctx *gin.Context) {
	postId := ctx.Param("post_id")

	// find post by id
	post := models.Post{}
	result := pc.DB.Preload("User", utils.SelectColumnDB("ID", "Name")).First(&post, "id = ?", postId)

	if result.Error != nil {
		switch result.Error.Error() {
		case "record not found":
			dtos.RespondWithError(ctx, http.StatusNotFound, "no post with the given id")
		default:
			dtos.RespondWithError(ctx, http.StatusBadGateway, result.Error.Error())
		}
		return
	}

	dtos.RespondWithJson(ctx, http.StatusOK, dtos.PostToPostResponse(&post))
}

func (pc *postController) GetAllPost(ctx *gin.Context) {
	// find all post
	posts := []models.Post{}
	result := pc.DB.Preload("User", utils.SelectColumnDB("ID", "Name")).Find(&posts)

	if result.Error != nil {
		dtos.RespondWithError(ctx, http.StatusBadGateway, result.Error.Error())
	}

	postsResponse := []*dtos.PostResponse{}
	for _, post := range posts {
		postsResponse = append(postsResponse, dtos.PostToPostResponse(&post))
	}

	dtos.RespondWithJson(ctx, http.StatusOK, postsResponse)
}

func (pc *postController) CreatePost(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	var payload *dtos.CreatePostInput

	// try to bind the request body to the payload struct
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		dtos.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	now := time.Now()

	newPost := models.Post{
		UserID:    currentUser.ID,
		Title:     payload.Title,
		Photo:     payload.Photo,
		Content:   payload.Content,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// save the data to database using gorm
	result := pc.DB.Create(&newPost)

	// check for any possible error
	if result.Error != nil {
		dtos.RespondWithError(ctx, http.StatusBadGateway, "something bad just happen")
		return
	}

	// give json response
	dtos.RespondWithJson(ctx, http.StatusCreated, dtos.PostToPostResponse(&newPost))
}

func (pc *postController) UpdatePost(ctx *gin.Context) {
	postId := ctx.Param("post_id")
	payload := dtos.UpdatePostInput{}

	// try to bind the request body to the payload struct
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		dtos.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// find the post
	var post models.Post
	result := pc.DB.First(&post, "id = ?", postId)
	if result.Error != nil {
		switch result.Error.Error() {
		case "record not found":
			dtos.RespondWithError(ctx, http.StatusNotFound, "no post with the given id")
		default:
			dtos.RespondWithError(ctx, http.StatusBadGateway, result.Error.Error())
		}
		return
	}

	// create the reflect value of the payload
	v := reflect.ValueOf(payload)

	// find value that want to be updated and update it
	tx := pc.DB.Begin()
	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldName := v.Type().Field(i).Name

		if !fieldValue.IsZero() {
			// update the post field
			fmt.Println(fieldValue)
			result := tx.Model(&post).Update(fieldName, fieldValue)
			if result.Error != nil {
				switch result.Error.Error() {
				case "record not found":
					dtos.RespondWithError(ctx, http.StatusNotFound, "no post with the given id")
				default:
					dtos.RespondWithError(ctx, http.StatusBadGateway, result.Error.Error())
				}
				return
			}
		}
	}
	tx.Commit()

	dtos.RespondWithJson(ctx, http.StatusOK, "successfully update post")
}

func (pc *postController) DeletePost(ctx *gin.Context) {
	postId := ctx.Param("post_id")

	result := pc.DB.Where("id = ?", postId).Delete(&models.Post{})
	if result.Error != nil {
		switch result.Error.Error() {
		case "record not found":
			dtos.RespondWithError(ctx, http.StatusNotFound, "no post with the given id")
		default:
			dtos.RespondWithError(ctx, http.StatusBadGateway, result.Error.Error())
		}
		return
	}

	dtos.RespondWithJson(ctx, http.StatusOK, "successfully delete post")
}
