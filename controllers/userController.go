package controllers

import "gorm.io/gorm"

type UserController struct {
	db *gorm.DB
}

type UserRegisterRequest struct {
	Age      uint   `json:"age"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type UserLoginRequest struct {
	Email    string `json:"email" valid:"required~email is required, email~Invalid format email"`
	Password string `json:"password" valid:"required~password is required, minstringlength(6)~password has to have minimum length of 6 characters"`
}

type UserUpdateRequest struct {
	Email    string `json:"email" valid:"email~Invalid format email"`
	Username string `json:"username"`
}
