package models

import (
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Comment struct {
	Id        uint       `gorm:"primaryKey" json:"id"`
	UserId    uint       `json:"user_id"`
	PhotoId   uint       `json:"photo_id"`
	Message   string     `gorm:"not null" json:"message" valid:"required~message is required"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`

	User  *User
	Photo *Photo
}

func (c *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(c)
	if errCreate != nil {
		return errCreate
	}
	return

}

func (c *Comment) BeforeUpdate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(c)
	if errCreate != nil {
		return errCreate
	}

	return
}
