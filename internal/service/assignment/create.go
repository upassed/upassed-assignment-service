package assignment

import (
	"context"
	"errors"
	"github.com/upassed/upassed-assignment-service/internal/async"
	"github.com/upassed/upassed-assignment-service/internal/handling"
	"github.com/upassed/upassed-assignment-service/internal/logging"
	"github.com/upassed/upassed-assignment-service/internal/middleware/common/auth"
	business "github.com/upassed/upassed-assignment-service/internal/service/model"
	"github.com/upassed/upassed-assignment-service/internal/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"log/slog"
)

var (
	ErrAssignmentCreateDeadlineExceeded = errors.New("assignment create deadline exceeded")
)

func (service *serviceImpl) Create(ctx context.Context, assignment *business.Assignment) (*business.AssignmentCreateResponse, error) {
	teacherUsername := ctx.Value(auth.UsernameKey).(string)

	spanContext, span := otel.Tracer(service.cfg.Tracing.AssignmentTracerName).Start(ctx, "assignmentService#Create")
	span.SetAttributes(
		attribute.String("formID", assignment.FormID.String()),
		attribute.String(auth.UsernameKey, teacherUsername),
	)
	defer span.End()

	log := logging.Wrap(service.log,
		logging.WithOp(service.Create),
		logging.WithCtx(ctx),
		logging.WithAny("formID", assignment.FormID),
	)

	log.Info("started creating assignment")
	timeout := service.cfg.GetEndpointExecutionTimeout()

	assignmentCreateResponse, err := async.ExecuteWithTimeout(spanContext, timeout, func(ctx context.Context) (*business.AssignmentCreateResponse, error) {
		log.Info("checking assignment duplicates")
		domainAssignment := ConvertToDomainAssignment(assignment)
		err := service.repository.CheckDuplicates(ctx, domainAssignment)
		if err != nil {
			log.Error("unable to check assignment duplicates", logging.Error(err))
			tracing.SetSpanError(span, err)
			return nil, err
		}

		log.Info("saving assignment data to the database")
		if err := service.repository.Save(ctx, domainAssignment); err != nil {
			log.Error("unable to save assignment data to the database", logging.Error(err))
			tracing.SetSpanError(span, err)
			return nil, err
		}

		return &business.AssignmentCreateResponse{
			CreatedAssignmentID: domainAssignment.ID,
		}, nil
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Error("assignment creating deadline exceeded")
			tracing.SetSpanError(span, err)
			return nil, handling.Wrap(ErrAssignmentCreateDeadlineExceeded, handling.WithCode(codes.DeadlineExceeded))
		}

		log.Error("error while creating assignment", logging.Error(err))
		tracing.SetSpanError(span, err)
		return nil, handling.Process(err)
	}

	log.Info("assignment successfully created", slog.Any("createdAssignmentID", assignmentCreateResponse.CreatedAssignmentID))
	return assignmentCreateResponse, nil
}
