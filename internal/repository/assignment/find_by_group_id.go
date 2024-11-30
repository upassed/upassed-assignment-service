package assignment

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/upassed/upassed-assignment-service/internal/handling"
	"github.com/upassed/upassed-assignment-service/internal/logging"
	domain "github.com/upassed/upassed-assignment-service/internal/repository/model"
	"github.com/upassed/upassed-assignment-service/internal/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
)

var (
	errFindingAssignmentsByGroupID = errors.New("error finding assignments by group id")
)

func (repository *repositoryImpl) FindByGroupID(ctx context.Context, groupID uuid.UUID) ([]*domain.Assignment, error) {
	spanContext, span := otel.Tracer(repository.cfg.Tracing.AssignmentTracerName).Start(ctx, "assignmentRepository#FindByGroupID")
	span.SetAttributes(
		attribute.String("groupID", groupID.String()),
	)
	defer span.End()

	log := logging.Wrap(repository.log,
		logging.WithOp(repository.FindByGroupID),
		logging.WithCtx(ctx),
		logging.WithAny("groupID", groupID),
	)

	log.Info("started finding assignments by group id in redis cache")
	cachedAssignments, err := repository.cache.GetByGroupID(spanContext, groupID)
	if err != nil || len(cachedAssignments) == 0 {
		log.Info("no cached assignments found, now going to the database")
	} else {
		log.Info("found assignments in the cache, not going to the database")
		return cachedAssignments, nil
	}

	log.Info("started finding assignments by group id in the database")
	foundAssignments := make([]*domain.Assignment, 0)
	findResult := repository.db.WithContext(ctx).Where("group_id = ?", groupID).Find(&foundAssignments)
	if err := findResult.Error; err != nil {
		log.Error("error while finding assignments by group id", logging.Error(err))
		tracing.SetSpanError(span, err)
		return nil, handling.New(errFindingAssignmentsByGroupID.Error(), codes.Internal)
	}

	log.Info("assignments by group id were successfully found in the database")

	if len(foundAssignments) != 0 {
		log.Info("saving assignments from the database to the cache")
		if err := repository.cache.SaveByGroupID(spanContext, foundAssignments); err != nil {
			log.Error("error while saving assignments from the database to the cache")
		}

		log.Info("assignments from the database were successfully saved to the cache")
	}

	return foundAssignments, nil
}
