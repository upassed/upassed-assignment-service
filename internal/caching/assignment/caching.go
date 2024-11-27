package assignment

import (
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/upassed/upassed-assignment-service/internal/config"
	"log/slog"
)

const (
	groupAssignmentsKeyFormat = "group-assignments:%s"
	formAssignmentsKeyFormat  = "form-assignments:%s"
)

var (
	errUnmarshallingAssignmentsDataFromJson = errors.New("unable to unmarshall assignments data from the cache from json format")
	errMarshallingAssignmentsDataToJson     = errors.New("unable to marshall assignments data to json format")
)

type RedisClient struct {
	cfg    *config.Config
	log    *slog.Logger
	client *redis.Client
}

func New(client *redis.Client, cfg *config.Config, log *slog.Logger) *RedisClient {
	return &RedisClient{
		cfg:    cfg,
		log:    log,
		client: client,
	}
}
