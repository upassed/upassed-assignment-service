package assignment

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/upassed/upassed-assignment-service/internal/logging"
	domain "github.com/upassed/upassed-assignment-service/internal/repository/model"
	"github.com/upassed/upassed-assignment-service/internal/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (client *RedisClient) SaveByGroupID(ctx context.Context, assignments []*domain.Assignment) error {
	groupID := assignments[0].GroupID

	_, span := otel.Tracer(client.cfg.Tracing.AssignmentTracerName).Start(ctx, "redisClient#SaveByGroupID")
	span.SetAttributes(attribute.String("groupID", groupID.String()))
	defer span.End()

	log := logging.Wrap(client.log,
		logging.WithOp(client.SaveByGroupID),
		logging.WithCtx(ctx),
		logging.WithAny("groupID", groupID),
	)

	log.Info("marshalling assignments data to json to save to the cache")
	jsonAssignmentsData, err := json.Marshal(assignments)
	if err != nil {
		log.Error("unable to marshall assignments data to json format")
		tracing.SetSpanError(span, err)
		return errMarshallingAssignmentsDataToJson
	}

	log.Info("started saving assignments data by groupID into the cache")
	if err := client.client.Set(ctx, fmt.Sprintf(groupAssignmentsKeyFormat, groupID), jsonAssignmentsData, client.cfg.GetRedisEntityTTL()).Err(); err != nil {
		log.Error("error while saving assignments data by groupID to the cache", logging.Error(err))
		tracing.SetSpanError(span, err)
		return errSavingAssignmentsByGroupIDIntoTheCache
	}

	log.Info("assignments were successfully saved by formID into the cache")
	return nil
}
