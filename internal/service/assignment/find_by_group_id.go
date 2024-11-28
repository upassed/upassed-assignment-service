package assignment

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/upassed/upassed-assignment-service/internal/async"
	"github.com/upassed/upassed-assignment-service/internal/handling"
	"github.com/upassed/upassed-assignment-service/internal/logging"
	"github.com/upassed/upassed-assignment-service/internal/middleware/common/auth"
	business "github.com/upassed/upassed-assignment-service/internal/service/model"
	"github.com/upassed/upassed-assignment-service/internal/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
)

var (
	ErrAssignmentFindByGroupIDDeadlineExceeded = errors.New("assignment find by group id deadline exceeded")
)

func (service *serviceImpl) FindByGroupID(ctx context.Context, groupID uuid.UUID) (*business.GroupAssignment, error) {
	teacherUsername := ctx.Value(auth.UsernameKey).(string)

	spanContext, span := otel.Tracer(service.cfg.Tracing.AssignmentTracerName).Start(ctx, "assignmentService#FindByGroupID")
	span.SetAttributes(
		attribute.String("groupID", groupID.String()),
		attribute.String(auth.UsernameKey, teacherUsername),
	)
	defer span.End()

	log := logging.Wrap(service.log,
		logging.WithOp(service.FindByGroupID),
		logging.WithCtx(ctx),
		logging.WithAny("groupID", groupID),
	)

	log.Info("started finding assignment by group id")
	timeout := service.cfg.GetEndpointExecutionTimeout()

	foundGroupAssignments, err := async.ExecuteWithTimeout(spanContext, timeout, func(ctx context.Context) (*business.GroupAssignment, error) {
		log.Info("searching assignment data by group id in the database")
		foundDomainAssignments, err := service.repository.FindByGroupID(ctx, groupID)
		if err != nil {
			log.Error("unable to find assignment data by group id in the database", logging.Error(err))
			tracing.SetSpanError(span, err)
			return nil, err
		}

		return ConvertToBusinessGroupAssignment(foundDomainAssignments), nil
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Error("assignment finding by group id deadline exceeded")
			tracing.SetSpanError(span, err)
			return nil, handling.Wrap(ErrAssignmentFindByGroupIDDeadlineExceeded, handling.WithCode(codes.DeadlineExceeded))
		}

		log.Error("error while finding assignment by group id", logging.Error(err))
		tracing.SetSpanError(span, err)
		return nil, handling.Process(err)
	}

	log.Info("assignment successfully found by group id")
	return foundGroupAssignments, nil
}
