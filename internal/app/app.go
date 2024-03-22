package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/Heatdog/VK-Go-PHP/internal/config"
	user_service "github.com/Heatdog/VK-Go-PHP/internal/service/user"
	user_handler "github.com/Heatdog/VK-Go-PHP/internal/transport/user"
	"github.com/Heatdog/VK-Go-PHP/pkg/client/postgre"
	"github.com/gorilla/mux"
)

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

	logger.Info("register user handler")
	userService := user_service.NewUserService(logger)
	userHandler := user_handler.NewUserHandler(logger, userService)
	userHandler.Register(router)

	host := fmt.Sprintf("%s:%s", cfg.Server.IP, cfg.Server.Port)
	logger.Info("listen tcp", slog.String("host", host))

	if err := http.ListenAndServe(host, router); err != nil {
		logger.Error(err.Error())
		panic(err)
	}
}
