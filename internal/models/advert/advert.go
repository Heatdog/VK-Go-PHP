package advert_model

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

const (
	min_price = 0
	max_price = 10000000

	max_height = 1080
	max_width  = 1920
)

func init() {
	govalidator.TagMap["image"] = govalidator.Validator(func(str string) bool {
		resp, err := http.Get(str)
		if err != nil {
			return false
		}
		defer resp.Body.Close()

		img, _, err := image.Decode(resp.Body)
		if err != nil {
			return false
		}

		height := img.Bounds().Dy()
		width := img.Bounds().Dx()

		if height < 0 || height > max_height {
			return false
		}
		if width < 0 || width > max_width {
			return false
		}
		return true
	})

	govalidator.TagMap["price"] = govalidator.Validator(func(str string) bool {
		val, err := strconv.Atoi(str)
		if err != nil {
			return false
		}
		if val < min_price || val > max_price {
			return false
		}
		return true
	})
}

type Advert struct {
	ID      uuid.UUID `json:"id"`
	Title   string    `json:"title" valid:"required,stringlength(3|250)"`
	Body    string    `json:"body" valid:"required,stringlength(3|1200)"`
	ImgAddr string    `json:"image" valid:"url,required"`
	Price   int       `json:"price" valid:"price,required"`
}

type AdvertWithOwner struct {
	Advert
	UserID    uuid.UUID `json:"user_id"`
	UserLogin string    `json:"user_login"`
	Own       bool      `json:"own"`
}

type AdvertInput struct {
	Title   string `json:"title" valid:"required,stringlength(3|250)"`
	Body    string `json:"body" valid:"required,stringlength(3|1200)"`
	ImgAddr string `json:"image" valid:"url,image,required"`
	Price   int    `json:"price" valid:"price,required"`
}

func MakeAdvert(in *AdvertInput, id uuid.UUID) *Advert {
	return &Advert{
		ID:      id,
		Title:   in.Title,
		Body:    in.Body,
		ImgAddr: in.ImgAddr,
		Price:   in.Price,
	}
}

type QueryParams struct {
	Sort     string
	SortDir  string
	MinPrice string `valid:"price"`
	MaxPrice string `valid:"price"`
}

func comapeStrInt(str1, str2 string) bool {
	val1, err := strconv.Atoi(str1)
	if err != nil {
		return false
	}

	val2, err := strconv.Atoi(str2)
	if err != nil {
		return false
	}

	if val1 > val2 {
		return false
	}
	return true
}

func fillQuery(params QueryParams) QueryParams {
	if params.Sort != "price" {
		params.Sort = "date_time"
	}
	if params.SortDir != "asc" {
		params.SortDir = "desc"
	}

	if params.MaxPrice == "" {
		params.MaxPrice = strconv.Itoa(max_price)
	}

	if params.MinPrice == "" {
		params.MinPrice = strconv.Itoa(min_price)
	}
	return params
}

func ValidQuery(params QueryParams) (QueryParams, error) {
	params = fillQuery(params)

	if _, err := govalidator.ValidateStruct(params); err != nil {
		return QueryParams{}, err
	}

	if !comapeStrInt(params.MinPrice, params.MaxPrice) {
		return QueryParams{}, fmt.Errorf("invalid price interval")
	}

	return params, nil
}
