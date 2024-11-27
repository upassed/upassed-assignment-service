package assignment

import (
	"context"
	"errors"
	"github.com/upassed/upassed-assignment-service/internal/handling"
	"github.com/upassed/upassed-assignment-service/internal/logging"
	domain "github.com/upassed/upassed-assignment-service/internal/repository/model"
	"github.com/upassed/upassed-assignment-service/internal/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"log/slog"
)

var (
	errSavingAssignments = errors.New("error while saving assignments into the database")
)

func (repository *repositoryImpl) Save(ctx context.Context, assignments []*domain.Assignment) error {
	spanContext, span := otel.Tracer(repository.cfg.Tracing.AssignmentTracerName).Start(ctx, "assignmentRepository#Save")
	span.SetAttributes(
		attribute.String("formID", assignments[0].FormID.String()),
	)
	defer span.End()

	log := logging.Wrap(repository.log,
		logging.WithOp(repository.CheckDuplicates),
		logging.WithCtx(ctx),
		logging.WithAny("formID", assignments[0].FormID),
	)

	log.Info("started saving assignments into the database")
	saveResult := repository.db.WithContext(ctx).CreateInBatches(assignments, 100)
	if err := saveResult.Error; err != nil {
		log.Error("error while saving assignments into the database", logging.Error(err))
		tracing.SetSpanError(span, err)
		return handling.New(errSavingAssignments.Error(), codes.Internal)
	}

	log.Info("saving assignments by formID into the redis cache")
	if err := repository.cache.SaveByFormID(spanContext, assignments); err != nil {
		log.Error("error while saving assignments by formID into the redis cache", logging.Error(err))
	} else {
		log.Info("assignments were successfully saved by formID into the cache")
	}

	for _, assignment := range assignments {
		log.Info("adding assignment to cached assignments by groupID", slog.Any("groupID", assignment.GroupID))
		if err := repository.cache.AddByGroupID(spanContext, assignment); err != nil {
			log.Error("error while adding assignment to cached assignments by groupID", logging.Error(err))
		}
	}

	log.Info("assignments were successfully inserted into the database")
	return nil
}
