package user_handler

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	user_model "github.com/Heatdog/VK-Go-PHP/internal/models/user"
	user_service "github.com/Heatdog/VK-Go-PHP/internal/service/user"
	"github.com/Heatdog/VK-Go-PHP/internal/transport"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
)

type userHandler struct {
	logger  *slog.Logger
	service user_service.UserService
}

func NewUserHandler(logger *slog.Logger, service user_service.UserService) transport.Handler {
	return &userHandler{
		logger:  logger,
		service: service,
	}
}

const (
	signUp = "register"
)

func (handler *userHandler) Register(router *mux.Router) {
	router.HandleFunc(signUp, handler.signUp).Methods(http.MethodPost)
}

// Регистрация в системе
// @Summary SignUp
// @Tags auth
// @Description Регистрациия в системе. Минимальная длина логина и пароля - 3 символа.
// @Description Логин должен быть уникальным.
// @ID sign-up
// @Accept json
// @Produce json
// @Param input body user_model.UserLogin true "user info"
// @Success 201 {object} transport.SignUpResponse Успешная регистрация
// @Failure 400 {object} transport.RespWriter Некооректные входные данные
// @Failure 500 {object} transport.RespWriter Внутренняя ошибка сервера
// @Router /register [post]
func (handler *userHandler) signUp(w http.ResponseWriter, r *http.Request) {
	handler.logger.Info("sign up handler")

	handler.logger.Debug("read request body")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		handler.logger.Warn(err.Error())
		transport.NewRespWriter(w, err.Error(), http.StatusBadRequest, handler.logger)
		return
	}
	defer r.Body.Close()

	handler.logger.Debug("request body", slog.String("body", string(body)))

	handler.logger.Debug("unmarshaling request body")
	var user user_model.UserLogin
	if err := json.Unmarshal(body, &user); err != nil {
		handler.logger.Warn(err.Error())
		transport.NewRespWriter(w, err.Error(), http.StatusBadRequest, handler.logger)
		return
	}

	handler.logger.Debug("validate user struct")
	_, err = govalidator.ValidateStruct(user)
	if err != nil {
		handler.logger.Warn(err.Error())
		transport.NewRespWriter(w, err.Error(), http.StatusBadRequest, handler.logger)
		return
	}

	handler.logger.Debug("user service")
	id, err := handler.service.SignUp(r.Context(), &user)
	if err != nil {
		handler.logger.Warn(err.Error())
		transport.NewRespWriter(w, err.Error(), http.StatusInternalServerError, handler.logger)
		return
	}

	handler.logger.Debug("marshal response", slog.String("id", id.String()), slog.String("login", user.Login))
	resp, err := json.Marshal(transport.SignUpResponse{
		ID:    id,
		Login: user.Login,
	})

	if err != nil {
		handler.logger.Warn(err.Error())
		transport.NewRespWriter(w, err.Error(), http.StatusInternalServerError, handler.logger)
		return
	}

	w.Write(resp)
	w.WriteHeader(http.StatusCreated)
	handler.logger.Info("successfull user registration", slog.String("id", id.String()),
		slog.String("login", user.Login))
}
