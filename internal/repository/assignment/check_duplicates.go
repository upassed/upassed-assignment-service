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
	"log/slog"
)

var (
	errCountingDuplicateAssignments = errors.New("error while counting assignments with formID and groupIDs in database")
)

func (repository *repositoryImpl) CheckDuplicates(ctx context.Context, assignments []*domain.Assignment) ([]*domain.Assignment, error) {
	_, span := otel.Tracer(repository.cfg.Tracing.AssignmentTracerName).Start(ctx, "assignmentRepository#CheckDuplicates")
	span.SetAttributes(
		attribute.String("formID", assignments[0].FormID.String()),
	)
	defer span.End()

	log := logging.Wrap(repository.log,
		logging.WithOp(repository.CheckDuplicates),
		logging.WithCtx(ctx),
		logging.WithAny("formID", assignments[0].FormID),
	)

	log.Info("started checking assignment duplicates")
	foundAssignments := make([]*domain.Assignment, 0)
	groupIDs := repository.mergeGroupIDs(assignments)

	searchResult := repository.db.WithContext(ctx).Model(&domain.Assignment{}).Where("form_id = ?", assignments[0].FormID).Where("group_id in (?)", groupIDs).Find(&foundAssignments)
	if err := searchResult.Error; err != nil {
		log.Error("error while counting assignments with formID and groupIDs in database", logging.Error(err))
		tracing.SetSpanError(span, err)
		return nil, handling.New(errCountingDuplicateAssignments.Error(), codes.Internal)
	}

	log.Info("assignment duplicates were searched in database", slog.Int("assignmentDuplicatesCount", len(foundAssignments)))
	return foundAssignments, nil
}

func (repository *repositoryImpl) mergeGroupIDs(assignments []*domain.Assignment) []uuid.UUID {
	groupIDs := make([]uuid.UUID, 0, len(assignments))
	for _, assignment := range assignments {
		groupIDs = append(groupIDs, assignment.GroupID)
	}

	return groupIDs
}
