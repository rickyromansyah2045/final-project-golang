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

type PhotoController struct {
	db *gorm.DB
}

type PhotoCreateRequest struct {
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoUrl string `json:"photo_url"`
}

type PhotoCreateResponse struct {
	Id        uint       `json:"id"`
	Title     string     `json:"title"`
	Caption   string     `json:"caption"`
	PhotoUrl  string     `json:"photo_url"`
	UserId    uint       `json:"user_id"`
	CreatedAt *time.Time `json:"created_at"`
}

type PhotoUpdateResponse struct {
	Id        uint       `json:"id"`
	Title     string     `json:"title"`
	Caption   string     `json:"caption"`
	PhotoUrl  string     `json:"photo_url"`
	UserId    uint       `json:"user_id"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type PhotoGetResponse struct {
	Id        uint       `json:"id"`
	Title     string     `json:"title"`
	Caption   string     `json:"caption"`
	PhotoUrl  string     `json:"photo_url"`
	UserId    uint       `json:"user_id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	User      UserDataResponse
}

type UserDataResponse struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

func NewPhotoController(db *gorm.DB) *PhotoController {
	return &PhotoController{
		db: db,
	}
}

func (p *PhotoController) Create(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	var photoReq PhotoCreateRequest

	err := ctx.ShouldBindJSON(&photoReq)
	if err != nil {
		helpers.BadRequestResponse(ctx, err)
		return
	}

	newPhoto := models.Photo{
		Title:    photoReq.Title,
		Caption:  photoReq.Caption,
		PhotoUrl: photoReq.PhotoUrl,
		UserId:   uint(userId.(float64)),
	}

	err = p.db.Create(&newPhoto).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helpers.NotFoundResponse(ctx, err)
			return
		}
		helpers.InternalServerJsonResponse(ctx, err)
		return
	}

	response := PhotoCreateResponse{
		Id:        newPhoto.Id,
		Title:     newPhoto.Title,
		Caption:   newPhoto.Caption,
		PhotoUrl:  newPhoto.PhotoUrl,
		UserId:    newPhoto.UserId,
		CreatedAt: newPhoto.CreatedAt,
	}

	helpers.WriteJsonResponse(ctx, http.StatusCreated, response)
}

func (p *PhotoController) Get(ctx *gin.Context) {
	var photos []models.Photo

	err := p.db.Preload("User").Find(&photos).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helpers.NotFoundResponse(ctx, err)
			return
		}
		helpers.BadRequestResponse(ctx, err)
		return
	}

	var response []PhotoGetResponse
	for _, photo := range photos {
		var userData UserDataResponse
		if photo.User != nil {
			userData = UserDataResponse{
				Username: photo.User.Username,
				Email:    photo.User.Email,
			}
		}
		response = append(response, PhotoGetResponse{
			Id:        photo.Id,
			Title:     photo.Title,
			Caption:   photo.Caption,
			PhotoUrl:  photo.PhotoUrl,
			UserId:    photo.UserId,
			CreatedAt: photo.CreatedAt,
			UpdatedAt: photo.UpdatedAt,
			User:      userData,
		})
	}

	helpers.WriteJsonResponse(ctx, http.StatusOK, response)
}

func (p *PhotoController) Update(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	photoId := ctx.Param("photoId")
	var photoReq PhotoCreateRequest
	var photo models.Photo

	err := ctx.ShouldBindJSON(&photoReq)
	if err != nil {
		helpers.BadRequestResponse(ctx, err)
		return
	}

	updatedPhoto := models.Photo{
		Title:    photoReq.Title,
		Caption:  photoReq.Caption,
		PhotoUrl: photoReq.PhotoUrl,
	}

	// Tambahin validasi Update
	_, errCreate := govalidator.ValidateStruct(&updatedPhoto)
	if errCreate != nil {
		helpers.BadRequestResponse(ctx, err)
		return
	}

	err = p.db.First(&photo, photoId).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helpers.NotFoundResponse(ctx, "data not found")
			return
		}
		helpers.InternalServerJsonResponse(ctx, err)
		return
	}

	if photo.UserId != uint(userId.(float64)) {
		helpers.UnauthorizeJsonResponse(ctx, "you're not allowed to update or edit this photo")
		return
	}

	err = p.db.Model(&photo).Updates(updatedPhoto).Error
	if err != nil {
		helpers.BadRequestResponse(ctx, err)
		return
	}

	response := PhotoUpdateResponse{
		Id:        photo.Id,
		Title:     photo.Title,
		Caption:   photo.Caption,
		PhotoUrl:  photo.PhotoUrl,
		UserId:    photo.UserId,
		UpdatedAt: photo.UpdatedAt,
	}

	helpers.WriteJsonResponse(ctx, http.StatusOK, response)
}

func (p *PhotoController) Delete(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	photoId := ctx.Param("photoId")
	var photo models.Photo

	err := p.db.First(&photo, photoId).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helpers.NotFoundResponse(ctx, "data not found")
			return
		}
		helpers.InternalServerJsonResponse(ctx, err)
		return
	}

	if photo.UserId != uint(userId.(float64)) {
		helpers.UnauthorizeJsonResponse(ctx, "you're not allowed to delete this photo")
		return
	}

	err = p.db.Delete(&photo).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helpers.NotFoundResponse(ctx, err.Error())
			return
		}
		helpers.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	helpers.WriteJsonResponse(ctx, http.StatusOK, gin.H{
		"message": "Your photo has been successfully deleted",
	})

}
