package controllers

import (
	"final-project-golang/helpers"
	"final-project-golang/models"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CommentController struct {
	db *gorm.DB
}

type CommentCreateRequest struct {
	Message string `json:"message"`
	PhotoId uint   `json:"photo_id"`
}

type CommentCreateResponse struct {
	Id        uint       `json:"id"`
	Message   string     `json:"message"`
	PhotoId   uint       `json:"photo_id"`
	UserId    uint       `json:"user_id"`
	CreatedAt *time.Time `json:"created_at"`
}

type CommentUpdateResponse struct {
	Id        uint       `json:"id"`
	Message   string     `json:"message"`
	PhotoId   uint       `json:"photo_id"`
	UserId    uint       `json:"user_id"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type CommentGetResponse struct {
	Id        uint       `json:"id"`
	Message   string     `json:"message"`
	PhotoId   uint       `json:"photo_id"`
	UserId    uint       `json:"user_id"`
	UpdatedAt *time.Time `json:"updated_at"`
	CreatedAt *time.Time `json:"created_at"`
	User      UserCommentResponse
	Photo     PhotoCommentResponse
}

type UserCommentResponse struct {
	Id       uint   `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type PhotoCommentResponse struct {
	Id       uint   `json:"id"`
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoUrl string `json:"photo_url"`
	UserId   uint   `json:"user_id"`
}

func NewCommentController(db *gorm.DB) *CommentController {
	return &CommentController{
		db: db,
	}
}

func (c *CommentController) Create(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	var commentReq CommentCreateRequest

	err := ctx.ShouldBindJSON(&commentReq)
	if err != nil {
		helpers.BadRequestResponse(ctx, err.Error())
		return
	}

	newComment := models.Comment{
		Message: commentReq.Message,
		PhotoId: commentReq.PhotoId,
		UserId:  uint(userId.(float64)),
	}

	_, errCreate := govalidator.ValidateStruct(&newComment)
	if errCreate != nil {
		helpers.BadRequestResponse(ctx, err.Error())
		return
	}

	err = c.db.Create(&newComment).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helpers.NotFoundResponse(ctx, err.Error())
			return
		}
		if err.Error() == `ERROR: insert or update on table "comments" violates foreign key constraint "fk_photos_comment" (SQLSTATE 23503)` {
			helpers.NotFoundResponse(ctx, "Photo not found")
			return
		}
		helpers.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	response := CommentCreateResponse{
		Id:        newComment.Id,
		Message:   newComment.Message,
		PhotoId:   newComment.PhotoId,
		UserId:    newComment.UserId,
		CreatedAt: newComment.CreatedAt,
	}

	helpers.WriteJsonResponse(ctx, http.StatusCreated, response)
}

func (c *CommentController) Get(ctx *gin.Context) {
	var comments []models.Comment

	err := c.db.Preload("User").Preload("Photo").Find(&comments).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helpers.NotFoundResponse(ctx, err.Error())
			return
		}
		helpers.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	var response []CommentGetResponse
	for _, comment := range comments {
		var userData UserCommentResponse
		if comment.User != nil {
			userData = UserCommentResponse{
				Id:       comment.User.Id,
				Username: comment.User.Username,
				Email:    comment.User.Email,
			}
		}
		var photoData PhotoCommentResponse
		if comment.Photo != nil {
			photoData = PhotoCommentResponse{
				Id:       comment.Photo.Id,
				Title:    comment.Photo.Title,
				Caption:  comment.Photo.Caption,
				PhotoUrl: comment.Photo.PhotoUrl,
				UserId:   comment.Photo.UserId,
			}
		}
		response = append(response, CommentGetResponse{
			Id:        comment.Id,
			Message:   comment.Message,
			PhotoId:   comment.PhotoId,
			UserId:    comment.UserId,
			UpdatedAt: comment.UpdatedAt,
			CreatedAt: comment.CreatedAt,
			User:      userData,
			Photo:     photoData,
		})
	}

	helpers.WriteJsonResponse(ctx, http.StatusOK, response)
}

func (c *CommentController) Update(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	commentId := ctx.Param("commentId")
	var commentReq CommentCreateRequest
	var comment models.Comment

	err := ctx.ShouldBindJSON(&commentReq)
	if err != nil {
		helpers.BadRequestResponse(ctx, err.Error())
		return
	}

	updateComment := models.Comment{
		Message: commentReq.Message,
	}

	err = c.db.First(&comment, commentId).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helpers.NotFoundResponse(ctx, "data not found")
			return
		}
		helpers.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	if comment.UserId != uint(userId.(float64)) {
		helpers.UnauthorizeJsonResponse(ctx, "you're not allowed to update or edit this comment")
		return
	}

	err = c.db.Model(&comment).Updates(updateComment).Error
	if err != nil {
		helpers.BadRequestResponse(ctx, err.Error())
		return
	}

	response := CommentUpdateResponse{
		Id:        comment.Id,
		Message:   comment.Message,
		PhotoId:   comment.PhotoId,
		UserId:    comment.UserId,
		UpdatedAt: comment.UpdatedAt,
	}

	helpers.WriteJsonResponse(ctx, http.StatusOK, response)
}

func (c *CommentController) Delete(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	commentId := ctx.Param("commentId")
	var comment models.Comment

	err := c.db.First(&comment, commentId).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helpers.NotFoundResponse(ctx, "data not found")
			return
		}
		helpers.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	if comment.UserId != uint(userId.(float64)) {
		helpers.BadRequestResponse(ctx, "you're not allowed to delete this comment")
		return
	}

	err = c.db.Delete(&comment).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helpers.NotFoundResponse(ctx, err.Error())
			return
		}
		helpers.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	helpers.WriteJsonResponse(ctx, http.StatusOK, gin.H{
		"message": "Your comment has been successfully deleted",
	})
}
