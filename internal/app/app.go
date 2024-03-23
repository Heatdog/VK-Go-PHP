package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/Heatdog/VK-Go-PHP/docs"

	"github.com/Heatdog/VK-Go-PHP/internal/config"
	user_postgre "github.com/Heatdog/VK-Go-PHP/internal/repository/user/postgre"
	token_service "github.com/Heatdog/VK-Go-PHP/internal/service/token"
	user_service "github.com/Heatdog/VK-Go-PHP/internal/service/user"
	user_handler "github.com/Heatdog/VK-Go-PHP/internal/transport/user"
	"github.com/Heatdog/VK-Go-PHP/pkg/client/postgre"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// swag init -g internal/app/app.go

// @title Маркетплейс
// @description API server for Маркетплейс

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
func App() {
	opt := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opt))
	slog.SetDefault(logger)

	ctx := context.Background()

	logger.Info("reading server config files")
	cfg := config.NewConfigStorage(logger)

	logger.Info("connecting to DataBase")
	dbClient, err := postgre.NewPostgreClient(ctx, cfg.Postgre)
	if err != nil {
		logger.Error("connection to PostgreSQL failed", slog.Any("error", err))
	}
	defer dbClient.Close()

	router := mux.NewRouter()

	logger.Info("reпister token service")
	tokenService := token_service.NewTokenService(logger, cfg.PasswordKey)

	logger.Info("register user handler")
	userRepo := user_postgre.NewUserPostgreRepository(dbClient, logger)
	userService := user_service.NewUserService(logger, userRepo, tokenService)
	userHandler := user_handler.NewUserHandler(logger, userService)
	userHandler.Register(router)

	logger.Info("adding swagger documentation")
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	host := fmt.Sprintf("%s:%s", cfg.Server.IP, cfg.Server.Port)
	logger.Info("listen tcp", slog.String("host", host))

	if err := http.ListenAndServe(host, router); err != nil {
		logger.Error(err.Error())
		panic(err)
	}
}
