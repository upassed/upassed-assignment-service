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

func (client *RedisClient) AddByGroupID(ctx context.Context, assignment *domain.Assignment) error {
	spanContext, span := otel.Tracer(client.cfg.Tracing.AssignmentTracerName).Start(ctx, "redisClient#AddByGroupID")
	span.SetAttributes(attribute.String("groupID", assignment.GroupID.String()))
	defer span.End()

	log := logging.Wrap(client.log,
		logging.WithOp(client.AddByGroupID),
		logging.WithCtx(ctx),
		logging.WithAny("groupID", assignment.GroupID),
	)

	existingGroupAssignments, err := client.GetByGroupID(spanContext, assignment.GroupID)
	if err != nil {
		log.Error("unable to get existing assignments data by groupID")
		tracing.SetSpanError(span, err)
		return errFetchingAssignmentsByGroupIDFromCache
	}

	newGroupAssignments := append(existingGroupAssignments, assignment)

	log.Info("marshalling assignments data to json to save to the cache")
	jsonAssignmentsData, err := json.Marshal(newGroupAssignments)
	if err != nil {
		log.Error("unable to marshall assignments data to json format")
		tracing.SetSpanError(span, err)
		return errMarshallingAssignmentsDataToJson
	}

	log.Info("started saving assignments data by groupID into the cache")
	if err := client.client.Set(ctx, fmt.Sprintf(groupAssignmentsKeyFormat, assignment.GroupID.String()), jsonAssignmentsData, client.cfg.GetRedisEntityTTL()).Err(); err != nil {
		log.Error("error while saving assignments data by groupID to the cache", logging.Error(err))
		tracing.SetSpanError(span, err)
		return errSavingAssignmentsByGroupIDIntoTheCache
	}

	log.Info("assignments were successfully saved by groupID into the cache")
	return nil
}
