package user_model

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `json:"id"`
	Login    string    `json:"login" valid:"required,length(3|50)"`
	Password string    `json:"password" valid:"required,length(3|50)"`
}

type UserLogin struct {
	Login    string `json:"login" valid:"required,length(3|50)"`
	Password string `json:"password" valid:"required,length(3|50)"`
}
