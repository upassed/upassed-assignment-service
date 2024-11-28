package app

import (
	"github.com/upassed/upassed-assignment-service/internal/caching"
	"github.com/upassed/upassed-assignment-service/internal/config"
	"github.com/upassed/upassed-assignment-service/internal/logging"
	"github.com/upassed/upassed-assignment-service/internal/messanging"
	assignmentRabbit "github.com/upassed/upassed-assignment-service/internal/messanging/assignment"
	"github.com/upassed/upassed-assignment-service/internal/middleware/common/auth"
	"github.com/upassed/upassed-assignment-service/internal/repository"
	assignmentRepo "github.com/upassed/upassed-assignment-service/internal/repository/assignment"
	"github.com/upassed/upassed-assignment-service/internal/server"
	assignmentSvc "github.com/upassed/upassed-assignment-service/internal/service/assignment"
	"log/slog"
)

type App struct {
	Server *server.AppServer
}

func New(config *config.Config, log *slog.Logger) (*App, error) {
	log = logging.Wrap(log, logging.WithOp(New))

	db, err := repository.OpenGormDbConnection(config, log)
	if err != nil {
		return nil, err
	}

	redis, err := caching.OpenRedisConnection(config, log)
	if err != nil {
		return nil, err
	}

	rabbit, err := messanging.OpenRabbitConnection(config, log)
	if err != nil {
		return nil, err
	}

	authClient, err := auth.NewClient(config, log)
	if err != nil {
		return nil, err
	}

	assignmentRepository := assignmentRepo.New(db, redis, config, log)
	assignmentService := assignmentSvc.New(config, log, assignmentRepository)
	assignmentRabbit.Initialize(authClient, assignmentService, rabbit, config, log)

	appServer := server.New(server.AppServerCreateParams{
		Config:            config,
		Log:               log,
		AuthClient:        authClient,
		AssignmentService: assignmentService,
	})

	log.Info("app successfully created")
	return &App{
		Server: appServer,
	}, nil
}
