package advert_handler

import (
	"log/slog"
	"net/http"

	"github.com/Heatdog/VK-Go-PHP/internal/transport"
	"github.com/gorilla/mux"
)

type advertHandler struct {
	logger *slog.Logger
}

func NewUserHandler(logger *slog.Logger) transport.Handler {
	return &advertHandler{
		logger: logger,
	}
}

const (
	add = "/advert/add"
)

func (handler *advertHandler) Register(router *mux.Router) {
	router.HandleFunc(add, handler.addAdvert).Methods(http.MethodPost)
}

// Добавление объявления
// @Summary AddAdvert
// @Tags advert
// @Description Добавление объявления в систему. Добавлять могут только авторизованные пользователи.
// @ID add-advert
// @Accept json
// @Produce json
// @Param input body user_model.UserLogin true "user info"
// @Success 201 {object} transport.SignUpResponse Успешная регистрация
// @Failure 400 {object} transport.RespWriter Некооректные входные данные
// @Failure 500 {object} transport.RespWriter Внутренняя ошибка сервера
// @Router /register [post]
func (handler *advertHandler) addAdvert(w http.ResponseWriter, r *http.Request) {

}
