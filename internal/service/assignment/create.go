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
	"log/slog"
)

var (
	ErrAssignmentCreateDeadlineExceeded = errors.New("assignment create deadline exceeded")
	ErrDuplicateAssignmentsFound        = errors.New("found duplicate assignments")
)

func (service *serviceImpl) Create(ctx context.Context, assignment *business.FormAssignment) (*business.AssignmentCreateResponse, error) {
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
		domainAssignments := ConvertToDomainAssignments(assignment)
		duplicateAssignments, err := service.repository.CheckDuplicates(ctx, domainAssignments)
		if err != nil {
			log.Error("unable to check assignment duplicates", logging.Error(err))
			tracing.SetSpanError(span, err)
			return nil, err
		}

		if len(duplicateAssignments) != 0 {
			log.Error("found duplicate assignments", slog.Int("assignmentDuplicatesCount", len(duplicateAssignments)))
			tracing.SetSpanError(span, ErrDuplicateAssignmentsFound)
			return nil, handling.New(ErrDuplicateAssignmentsFound.Error(), codes.AlreadyExists)
		}

		log.Info("saving assignment data to the database")
		if err := service.repository.Save(ctx, domainAssignments); err != nil {
			log.Error("unable to save assignment data to the database", logging.Error(err))
			tracing.SetSpanError(span, err)
			return nil, err
		}

		createdAssignmentIDs := make([]uuid.UUID, 0, len(domainAssignments))
		for _, domainAssignment := range domainAssignments {
			createdAssignmentIDs = append(createdAssignmentIDs, domainAssignment.ID)
		}

		return &business.AssignmentCreateResponse{
			CreatedAssignmentIDs: createdAssignmentIDs,
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

	log.Info("assignment successfully created", slog.Any("createdAssignmentIDs", assignmentCreateResponse.CreatedAssignmentIDs))
	return assignmentCreateResponse, nil
}
