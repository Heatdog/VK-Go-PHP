package user_model

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `json:"id"`
	Login    string    `json:"login" valid:"required,range(3|250)"`
	Password string    `json:"password" valid:"required,range(3|250)"`
}

type UserLogin struct {
	Login    string `json:"login" valid:"required,range(3|250)"`
	Password string `json:"password" valid:"required,range(3|250)"`
}
