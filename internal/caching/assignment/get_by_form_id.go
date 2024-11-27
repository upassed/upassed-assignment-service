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
	errFetchingAssignmentsByFormIDFromCache = errors.New("unable to get assignments by formID from the cache")
)

func (client *RedisClient) GetByFormID(ctx context.Context, formID uuid.UUID) ([]*domain.Assignment, error) {
	_, span := otel.Tracer(client.cfg.Tracing.AssignmentTracerName).Start(ctx, "redisClient#GetByFormID")
	span.SetAttributes(attribute.String("formID", formID.String()))
	defer span.End()

	log := logging.Wrap(client.log,
		logging.WithOp(client.GetByFormID),
		logging.WithCtx(ctx),
		logging.WithAny("formID", formID),
	)

	log.Info("started getting assignments data by formID from cache")
	formAssignmentsData, err := client.client.Get(ctx, fmt.Sprintf(formAssignmentsKeyFormat, formID.String())).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Info("assignments by formID were not found in cache")
			return make([]*domain.Assignment, 0), nil
		}

		log.Error("error while fetching assignments by formID from cache", logging.Error(err))
		tracing.SetSpanError(span, err)
		return nil, errFetchingAssignmentsByFormIDFromCache
	}

	log.Info("assignments by formID were found in cache, unmarshalling from json")
	assignments := make([]*domain.Assignment, 0)
	if err := json.Unmarshal([]byte(formAssignmentsData), &assignments); err != nil {
		log.Error("error while unmarshalling assignments data to json", logging.Error(err))
		tracing.SetSpanError(span, err)
		return nil, errUnmarshallingAssignmentsDataFromJson
	}

	log.Info("assignments was successfully found and unmarshalled")
	return assignments, nil
}
