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

type SocialController struct {
	db *gorm.DB
}

type SocialCreateRequest struct {
	Name           string `json:"name"`
	SocialMediaUrl string `json:"social_media_url"`
}

type SocialCreateResponse struct {
	Id             uint       `json:"id"`
	Name           string     `json:"name"`
	SocialMediaUrl string     `json:"social_media_url"`
	UserId         uint       `json:"user_id"`
	CreatedAt      *time.Time `json:"created_at"`
}

type SocialUpdateResponse struct {
	Id             uint       `json:"id"`
	Name           string     `json:"name"`
	SocialMediaUrl string     `json:"social_media_url"`
	UserId         uint       `json:"user_id"`
	UpdatedAt      *time.Time `json:"updated_at"`
}

type SocialGetResponse struct {
	SocialMedias []SocialData `json:"social_medias"`
}

type SocialData struct {
	Id             uint       `json:"id"`
	Name           string     `json:"name"`
	SocialMediaUrl string     `json:"social_media_url"`
	UserId         uint       `json:"user_id"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
	User           UserSocialResponse
}

type UserSocialResponse struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
}

func NewSocialController(db *gorm.DB) *SocialController {
	return &SocialController{
		db: db,
	}
}

func (s *SocialController) Create(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	var socialReq SocialCreateRequest

	err := ctx.ShouldBindJSON(&socialReq)
	if err != nil {
		helpers.BadRequestResponse(ctx, err.Error())
		return
	}

	newSocial := models.Social{
		Name:           socialReq.Name,
		SocialMediaUrl: socialReq.SocialMediaUrl,
		UserId:         uint(userId.(float64)),
	}

	_, errCreate := govalidator.ValidateStruct(&newSocial)
	if errCreate != nil {
		helpers.BadRequestResponse(ctx, errCreate.Error())
		return
	}

	err = s.db.Create(&newSocial).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helpers.NotFoundResponse(ctx, err.Error())
			return
		}
		helpers.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	response := SocialCreateResponse{
		Id:             newSocial.Id,
		Name:           newSocial.Name,
		SocialMediaUrl: newSocial.SocialMediaUrl,
		UserId:         newSocial.UserId,
		CreatedAt:      newSocial.CreatedAt,
	}

	helpers.WriteJsonResponse(ctx, http.StatusCreated, response)
}

func (s *SocialController) Get(ctx *gin.Context) {
	var socials []models.Social

	err := s.db.Preload("User").Find(&socials).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helpers.NotFoundResponse(ctx, err.Error())
			return
		}
		helpers.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	var response SocialGetResponse
	for _, social := range socials {
		var userData UserSocialResponse
		if social.User != nil {
			userData = UserSocialResponse{
				Id:       social.User.Id,
				Username: social.User.Username,
			}
		}

		socialMediasResponse := SocialData{
			Id:             social.Id,
			Name:           social.Name,
			SocialMediaUrl: social.SocialMediaUrl,
			UserId:         social.UserId,
			CreatedAt:      social.CreatedAt,
			UpdatedAt:      social.UpdatedAt,
			User:           userData,
		}

		response.SocialMedias = append(response.SocialMedias, socialMediasResponse)
	}

	if len(response.SocialMedias) == 0 {
		response.SocialMedias = make([]SocialData, 0)
	}

	helpers.WriteJsonResponse(ctx, http.StatusOK, response)
}

func (s *SocialController) Update(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	socialMediaId := ctx.Param("socialMediaId")
	var socialReq SocialCreateRequest
	var social models.Social

	err := ctx.ShouldBindJSON(&socialReq)
	if err != nil {
		helpers.BadRequestResponse(ctx, err.Error())
		return
	}

	updatedSocial := models.Social{
		Name:           socialReq.Name,
		SocialMediaUrl: socialReq.SocialMediaUrl,
		UserId:         uint(userId.(float64)),
	}

	err = s.db.First(&social, socialMediaId).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helpers.NotFoundResponse(ctx, "data not found")
			return
		}
		helpers.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	if social.UserId != uint(userId.(float64)) {
		helpers.UnauthorizeJsonResponse(ctx, "you're not allowed to update or edit this social media")
		return
	}

	err = s.db.Model(&social).Updates(updatedSocial).Error
	if err != nil {
		helpers.BadRequestResponse(ctx, err.Error())
		return
	}

	response := SocialUpdateResponse{
		Id:             social.Id,
		Name:           social.Name,
		SocialMediaUrl: social.SocialMediaUrl,
		UserId:         social.UserId,
		UpdatedAt:      social.UpdatedAt,
	}

	helpers.WriteJsonResponse(ctx, http.StatusOK, response)
}

func (s *SocialController) Delete(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	socialId := ctx.Param("socialMediaId")
	var social models.Social

	err := s.db.First(&social, socialId).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helpers.NotFoundResponse(ctx, "data not found")
			return
		}
		helpers.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	if social.UserId != uint(userId.(float64)) {
		helpers.UnauthorizeJsonResponse(ctx, "you're not allowed to delete this social media")
		return
	}

	err = s.db.Delete(&social).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helpers.NotFoundResponse(ctx, err.Error())
			return
		}
		helpers.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	helpers.WriteJsonResponse(ctx, http.StatusOK, gin.H{
		"message": "Your social media has been successfully deleted",
	})

}
