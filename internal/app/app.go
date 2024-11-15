package app

import (
	"github.com/upassed/upassed-assignment-service/internal/config"
	"github.com/upassed/upassed-assignment-service/internal/logging"
	"github.com/upassed/upassed-assignment-service/internal/messanging"
	"github.com/upassed/upassed-assignment-service/internal/middleware/common/auth"
	"github.com/upassed/upassed-assignment-service/internal/repository"
	"github.com/upassed/upassed-assignment-service/internal/server"
	"log/slog"
)

type App struct {
	Server *server.AppServer
}

func New(config *config.Config, log *slog.Logger) (*App, error) {
	log = logging.Wrap(log, logging.WithOp(New))

	_, err := repository.OpenGormDbConnection(config, log)
	if err != nil {
		return nil, err
	}

	_, err = messanging.OpenRabbitConnection(config, log)
	if err != nil {
		return nil, err
	}

	authClient, err := auth.NewClient(config, log)
	if err != nil {
		return nil, err
	}

	appServer, err := server.New(server.AppServerCreateParams{
		Config:     config,
		Log:        log,
		AuthClient: authClient,
	})

	if err != nil {
		log.Error("unable to create new grpc server", logging.Error(err))
		return nil, err
	}

	log.Info("app successfully created")
	return &App{
		Server: appServer,
	}, nil
}
