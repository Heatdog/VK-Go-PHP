package advert_handler

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	advert_model "github.com/Heatdog/VK-Go-PHP/internal/models/advert"
	advert_service "github.com/Heatdog/VK-Go-PHP/internal/service/advert"
	"github.com/Heatdog/VK-Go-PHP/internal/transport"
	middleware_transport "github.com/Heatdog/VK-Go-PHP/internal/transport/middleware"
	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type advertHandler struct {
	logger        *slog.Logger
	advertService advert_service.AdvertService
	middleware    *middleware_transport.Middleware
}

func NewUserHandler(logger *slog.Logger, adverService advert_service.AdvertService,
	middleware *middleware_transport.Middleware) transport.Handler {
	return &advertHandler{
		logger:        logger,
		advertService: adverService,
		middleware:    middleware,
	}
}

const (
	add = "/advert/add"
	get = "/advert/get"
)

func (handler *advertHandler) Register(router *mux.Router) {
	router.HandleFunc(add, handler.middleware.Auth(true, handler.addAdvert)).Methods(http.MethodPost)
	router.HandleFunc(get, handler.middleware.Auth(false, handler.getAdverts)).Methods(http.MethodGet)
}

// Добавление объявления
// @Summary AddAdvert
// @Security ApiKeyAuth
// @Tags advert
// @Description Добавление объявления в систему. Добавлять объявления могут только авторизованные пользователи.
// @Description Ограничение на загловок - от 3 до 250 символов; на текст объявления - от 3 до 1200 символов;
// @Description Формат изображения - jpg и png. Размер изображения - 1080 в длину и 1920 в ширину.
// @Description Ограничение цены - от 0 до 10 000 000
// @ID add-advert
// @Accept json
// @Produce json
// @Param input body advert_model.AdvertInput true "поля объявления"
// @Success 201 {object} advert_model.Advert Успешное добавление объявления в систему
// @Failure 400 {object} transport.RespWriter Некооректные входные данные
// @Failure 401 {object} transport.RespWriter Пользователь неавториован
// @Failure 500 {object} transport.RespWriter Внутренняя ошибка сервера
// @Router /advert/add [post]
func (handler *advertHandler) addAdvert(w http.ResponseWriter, r *http.Request) {
	handler.logger.Info("add advert handler")

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
	var advert advert_model.AdvertInput

	if err := json.Unmarshal(body, &advert); err != nil {
		handler.logger.Warn(err.Error())
		transport.NewRespWriter(w, err.Error(), http.StatusBadRequest, handler.logger)
		return
	}

	handler.logger.Debug("validate advert struct")
	_, err = govalidator.ValidateStruct(advert)
	if err != nil {
		handler.logger.Warn(err.Error())
		transport.NewRespWriter(w, err.Error(), http.StatusBadRequest, handler.logger)
		return
	}

	handler.logger.Debug("user id", r.Context().Value("user_id"))
	userID := r.Context().Value("user_id").(uuid.UUID)

	handler.logger.Debug("advert service")
	id, err := handler.advertService.AddAdvert(r.Context(), &advert, userID)

	if err != nil {
		handler.logger.Warn(err.Error())
		transport.NewRespWriter(w, err.Error(), http.StatusInternalServerError, handler.logger)
		return
	}

	advertResp := advert_model.MakeAdvert(&advert, id)
	handler.logger.Debug("marshal response", slog.Any("advert response", advertResp))
	resp, err := json.Marshal(advertResp)

	if err != nil {
		handler.logger.Warn(err.Error())
		transport.NewRespWriter(w, err.Error(), http.StatusInternalServerError, handler.logger)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
	handler.logger.Info("successfull advert add", slog.Any("advert", advertResp))
}

// Получение объявлений
// @Summary getAdverts
// @Tags advert
// @Security ApiKeyAuth
// @Description Получение списка объявлений. Возможность сортировки по дате и цене,
// @Description также можно задать направление сортировки. Возможность фильтрации по цене с мин и макс значениями.
// @Description Сортировка задается параметрами URL: order и dir. Если order=price, то сортировка будет по цене.
// @Description Иначе - по дате добавлени. Если dir=asc, то сортировка будет по возрастанию. Иначе - по убыванию.
// @Description Параметры min и max - ограничения на цену. Проверяется ограничения на то, что min <= max
// @Description и не выходит за пределы ограничений по цене.
// @Description Создаются страницы по 10 объявлений.
// @ID get-adverts
// @Accept json
// @Produce json
// @Param order query string false "type of order"
// @Param dir query string false "asc or desc"
// @Param min query string false "min price"
// @Param max query string false "max price"
// @Success 200 {object} [][]advert_model.AdvertWithOwner Список объявлений
// @Failure 400 {object} transport.RespWriter Некооректные входные данные
// @Failure 500 {object} transport.RespWriter Внутренняя ошибка сервера
// @Router /advert/get [get]
func (handler *advertHandler) getAdverts(w http.ResponseWriter, r *http.Request) {
	handler.logger.Info("get advert handler")

	values := r.URL.Query()
	queryParams, err := advert_model.ValidQuery(advert_model.QueryParams{
		Sort:     values.Get("order"),
		SortDir:  values.Get("dir"),
		MinPrice: values.Get("min"),
		MaxPrice: values.Get("max"),
	})

	if err != nil {
		handler.logger.Warn(err.Error())
		transport.NewRespWriter(w, err.Error(), http.StatusBadRequest, handler.logger)
		return
	}

	userID := uuid.Nil
	u_id := r.Context().Value("user_id")
	if u_id != nil {
		handler.logger.Debug("user id", u_id)
		userID = u_id.(uuid.UUID)
	}

	handler.logger.Debug("get adverts", slog.Any("query", queryParams))
	list, err := handler.advertService.GetAdverts(r.Context(), queryParams, userID)

	if err != nil {
		handler.logger.Warn(err.Error())
		transport.NewRespWriter(w, err.Error(), http.StatusInternalServerError, handler.logger)
		return
	}

	handler.logger.Debug("marshal response", slog.Any("advert response", list))
	resp, err := json.Marshal(list)

	if err != nil {
		handler.logger.Warn(err.Error())
		transport.NewRespWriter(w, err.Error(), http.StatusInternalServerError, handler.logger)
		return
	}

	w.Write(resp)
	w.WriteHeader(http.StatusOK)
	handler.logger.Info("successfull adverts get", slog.Any("advert", resp))
}
