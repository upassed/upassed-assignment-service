package assignment

import (
	"context"
	"github.com/google/uuid"
	"github.com/upassed/upassed-assignment-service/internal/config"
	business "github.com/upassed/upassed-assignment-service/internal/service/model"
	"github.com/upassed/upassed-assignment-service/pkg/client"
	"google.golang.org/grpc"
)

type assignmentServerAPI struct {
	client.UnimplementedAssignmentServer
	cfg     *config.Config
	service assignmentService
}

type assignmentService interface {
	FindByFormID(ctx context.Context, formID uuid.UUID) (*business.FormAssignment, error)
	FindByGroupID(ctx context.Context, groupID uuid.UUID) (*business.GroupAssignment, error)
}

func Register(gRPC *grpc.Server, cfg *config.Config, service assignmentService) {
	client.RegisterAssignmentServer(gRPC, &assignmentServerAPI{
		cfg:     cfg,
		service: service,
	})
}
