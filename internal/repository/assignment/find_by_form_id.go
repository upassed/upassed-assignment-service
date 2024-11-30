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
	errFindingAssignmentsByFormID = errors.New("error finding assignments by form id")
)

func (repository *repositoryImpl) FindByFormID(ctx context.Context, formID uuid.UUID) ([]*domain.Assignment, error) {
	spanContext, span := otel.Tracer(repository.cfg.Tracing.AssignmentTracerName).Start(ctx, "assignmentRepository#FindByFormID")
	span.SetAttributes(
		attribute.String("formID", formID.String()),
	)
	defer span.End()

	log := logging.Wrap(repository.log,
		logging.WithOp(repository.FindByFormID),
		logging.WithCtx(ctx),
		logging.WithAny("formID", formID),
	)

	log.Info("started finding assignments by form id in redis cache")
	cachedAssignments, err := repository.cache.GetByFormID(spanContext, formID)
	if err != nil || len(cachedAssignments) == 0 {
		log.Info("no cached assignments found, now going to the database")
	} else {
		log.Info("found assignments in the cache, not going to the database")
		return cachedAssignments, nil
	}

	log.Info("started finding assignments by form id in the database")
	foundAssignments := make([]*domain.Assignment, 0)
	findResult := repository.db.WithContext(ctx).Where("form_id = ?", formID).Find(&foundAssignments)
	if err := findResult.Error; err != nil {
		log.Error("error while finding assignments by form id", logging.Error(err))
		tracing.SetSpanError(span, err)
		return nil, handling.New(errFindingAssignmentsByFormID.Error(), codes.Internal)
	}

	log.Info("assignments by form id were successfully found in the database")

	if len(foundAssignments) != 0 {
		log.Info("saving assignments from the database to the cache")
		if err := repository.cache.SaveByFormID(spanContext, foundAssignments); err != nil {
			log.Error("error while saving assignments from the database to the cache")
		}

		log.Info("assignments from the database were successfully saved to the cache")
	}

	return foundAssignments, nil
}
