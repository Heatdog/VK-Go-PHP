package advert_model

import "github.com/google/uuid"

type Advert struct {
	ID      uuid.UUID `json:"id"`
	Title   string    `json:"title" valid:"required,length(3|250)"`
	Body    string    `json:"body" valid:"required,length(3|1200)"`
	ImgAddr string    `json:"image" valid:"url,required"`
	Price   int       `json:"price" valid:"required,range(0|)"`
}
