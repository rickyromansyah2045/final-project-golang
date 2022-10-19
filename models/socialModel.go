package models

import "time"

type Sosial struct {
	Id             uint   `gorm:"primaryKey" json:"id"`
	Name           string `gorm:"not null;type:varchar(100)" json:"name" valid:"required~name is required"`
	SocialMediaUrl string `gorm:"not null" json:"sosial_media_url" valid:"required~sosial_media_url is required"`
	UserId         uint   `json:"user_id"`
	CreatedAt      *time.Time
	UpdatedAt      *time.Time
}
