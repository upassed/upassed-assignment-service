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
	ErrAssignmentFindByFormIDDeadlineExceeded = errors.New("assignment find by form id deadline exceeded")
)

func (service *serviceImpl) FindByFormID(ctx context.Context, formID uuid.UUID) (*business.FormAssignment, error) {
	teacherUsername := ctx.Value(auth.UsernameKey).(string)

	spanContext, span := otel.Tracer(service.cfg.Tracing.AssignmentTracerName).Start(ctx, "assignmentService#FindByFormID")
	span.SetAttributes(
		attribute.String("formID", formID.String()),
		attribute.String(auth.UsernameKey, teacherUsername),
	)
	defer span.End()

	log := logging.Wrap(service.log,
		logging.WithOp(service.FindByFormID),
		logging.WithCtx(ctx),
		logging.WithAny("formID", formID),
	)

	log.Info("started finding assignment by form id")
	timeout := service.cfg.GetEndpointExecutionTimeout()

	foundFormAssignments, err := async.ExecuteWithTimeout(spanContext, timeout, func(ctx context.Context) (*business.FormAssignment, error) {
		log.Info("searching assignment data by form id in the database")
		foundDomainAssignments, err := service.repository.FindByFormID(ctx, formID)
		if err != nil {
			log.Error("unable to find assignment data by form id in the database", logging.Error(err))
			tracing.SetSpanError(span, err)
			return nil, err
		}

		return ConvertToBusinessFormAssignment(foundDomainAssignments), nil
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Error("assignment finding by form id deadline exceeded")
			tracing.SetSpanError(span, err)
			return nil, handling.Wrap(ErrAssignmentFindByFormIDDeadlineExceeded, handling.WithCode(codes.DeadlineExceeded))
		}

		log.Error("error while finding assignment by form id", logging.Error(err))
		tracing.SetSpanError(span, err)
		return nil, handling.Process(err)
	}

	log.Info("assignment successfully found by form id")
	return foundFormAssignments, nil
}
