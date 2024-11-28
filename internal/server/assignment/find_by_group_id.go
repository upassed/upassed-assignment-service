package assignment

import (
	"context"
	"github.com/google/uuid"
	"github.com/upassed/upassed-assignment-service/internal/handling"
	requestid "github.com/upassed/upassed-assignment-service/internal/middleware/common/request_id"
	"github.com/upassed/upassed-assignment-service/internal/tracing"
	"github.com/upassed/upassed-assignment-service/pkg/client"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
)

func (server *assignmentServerAPI) FindByGroupID(ctx context.Context, request *client.AssignmentFindByGroupIDRequest) (*client.AssignmentFindByGroupIDResponse, error) {
	spanContext, span := otel.Tracer(server.cfg.Tracing.AssignmentTracerName).Start(ctx, "assignment#FindByGroupID")
	span.SetAttributes(
		attribute.String(string(requestid.ContextKey), requestid.GetRequestIDFromContext(ctx)),
		attribute.String("groupID", request.GetGroupId()),
	)
	defer span.End()

	if err := request.Validate(); err != nil {
		tracing.SetSpanError(span, err)
		return nil, handling.Wrap(err, handling.WithCode(codes.InvalidArgument))
	}

	response, err := server.service.FindByGroupID(spanContext, uuid.MustParse(request.GetGroupId()))
	if err != nil {
		tracing.SetSpanError(span, err)
		return nil, err
	}

	return ConvertToFindByGroupIDResponse(response), nil
}
