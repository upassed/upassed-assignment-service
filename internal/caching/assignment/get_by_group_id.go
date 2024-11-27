package assignment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/upassed/upassed-assignment-service/internal/logging"
	domain "github.com/upassed/upassed-assignment-service/internal/repository/model"
	"github.com/upassed/upassed-assignment-service/internal/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var (
	errFetchingAssignmentsByGroupIDFromCache = errors.New("unable to get assignments by groupID from the cache")
)

func (client *RedisClient) GetByGroupID(ctx context.Context, groupID uuid.UUID) ([]*domain.Assignment, error) {
	_, span := otel.Tracer(client.cfg.Tracing.AssignmentTracerName).Start(ctx, "redisClient#GetByGroupID")
	span.SetAttributes(attribute.String("groupID", groupID.String()))
	defer span.End()

	log := logging.Wrap(client.log,
		logging.WithOp(client.GetByGroupID),
		logging.WithCtx(ctx),
		logging.WithAny("groupID", groupID),
	)

	log.Info("started getting assignments data by groupID from cache")
	groupAssignmentsData, err := client.client.Get(ctx, fmt.Sprintf(groupAssignmentsKeyFormat, groupID.String())).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Error("assignments by groupID were not found in cache")
			return make([]*domain.Assignment, 0), nil
		}

		log.Error("error while fetching assignments by groupID from cache", logging.Error(err))
		tracing.SetSpanError(span, err)
		return nil, errFetchingAssignmentsByGroupIDFromCache
	}

	log.Info("assignments by groupID were found in cache, unmarshalling from json")
	assignments := make([]*domain.Assignment, 0)
	if err := json.Unmarshal([]byte(groupAssignmentsData), &assignments); err != nil {
		log.Error("error while unmarshalling assignments data to json", logging.Error(err))
		tracing.SetSpanError(span, err)
		return nil, errUnmarshallingAssignmentsDataFromJson
	}

	log.Info("assignments was successfully found and unmarshalled")
	return assignments, nil
}
