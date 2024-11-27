package assignment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/upassed/upassed-assignment-service/internal/logging"
	domain "github.com/upassed/upassed-assignment-service/internal/repository/model"
	"github.com/upassed/upassed-assignment-service/internal/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var (
	errSavingAssignmentsByFormIDIntoTheCache = errors.New("unable to save assignments data by form id into the cache")
)

func (client *RedisClient) SaveByFormID(ctx context.Context, assignments []*domain.Assignment) error {
	formID := assignments[0].FormID.String()

	_, span := otel.Tracer(client.cfg.Tracing.AssignmentTracerName).Start(ctx, "redisClient#SaveByFormID")
	span.SetAttributes(attribute.String("formID", formID))
	defer span.End()

	log := logging.Wrap(client.log,
		logging.WithOp(client.SaveByFormID),
		logging.WithCtx(ctx),
		logging.WithAny("formID", formID),
	)

	log.Info("marshalling assignments data to json to save to the cache")
	jsonAssignmentsData, err := json.Marshal(assignments)
	if err != nil {
		log.Error("unable to marshall assignments data to json format")
		tracing.SetSpanError(span, err)
		return errMarshallingAssignmentsDataToJson
	}

	log.Info("started saving assignments data by formID into the cache")
	if err := client.client.Set(ctx, fmt.Sprintf(formAssignmentsKeyFormat, formID), jsonAssignmentsData, client.cfg.GetRedisEntityTTL()).Err(); err != nil {
		log.Error("error while saving assignments data by formID to the cache", logging.Error(err))
		tracing.SetSpanError(span, err)
		return errSavingAssignmentsByFormIDIntoTheCache
	}

	log.Info("assignments were successfully saved by formID into the cache")
	return nil
}
