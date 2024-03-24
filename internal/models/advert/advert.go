package advert_model

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
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

		if height < 0 || height > 1920 {
			return false
		}
		if width < 0 || width > 1080 {
			return false
		}
		return true
	})
	govalidator.TagMap["gte"] = govalidator.Validator(func(str string) bool {
		val, err := strconv.Atoi(str)
		if err != nil {
			return false
		}
		if val < 0 {
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
	Price   int       `json:"price" valid:"required,gte"`
	UserID  uuid.UUID `json:"user_id"`
	Own     bool      `json:"own"`
}

type AdvertInput struct {
	Title   string `json:"title" valid:"required,stringlength(3|250)"`
	Body    string `json:"body" valid:"required,stringlength(3|1200)"`
	ImgAddr string `json:"image" valid:"url,image,required"`
	Price   int    `json:"price" valid:"required,gte"`
}

func MakeAdvert(in *AdvertInput, id uuid.UUID, userID uuid.UUID) *Advert {
	return &Advert{
		ID:      id,
		Title:   in.Title,
		Body:    in.Body,
		ImgAddr: in.ImgAddr,
		Price:   in.Price,
		UserID:  userID,
		Own:     true,
	}
}

type QueryParams struct {
	Sort     string
	SortDir  string
	MinPrice string `valid:"int"`
	MaxPrice string `valid:"int"`
}

func ValidQuery(params QueryParams) (QueryParams, error) {
	if params.Sort != "price" {
		params.Sort = "date"
	}
	if params.SortDir != "asc" {
		params.SortDir = "desc"
	}

	if params.MaxPrice == "" {
		params.MaxPrice = strconv.Itoa(math.MaxInt)
	}

	if params.MinPrice == "" {
		params.MinPrice = strconv.Itoa(0)
	}

	if _, err := govalidator.ValidateStruct(params); err != nil {
		return QueryParams{}, err
	}

	minPrice, err := strconv.Atoi(params.MinPrice)
	if err != nil {
		return QueryParams{}, err
	}

	maxPrice, err := strconv.Atoi(params.MaxPrice)
	if err != nil {
		return QueryParams{}, err
	}

	if minPrice < 0 || maxPrice < 0 || minPrice < maxPrice {
		return QueryParams{}, fmt.Errorf("invalid price interval")
	}

	return params, nil
}
