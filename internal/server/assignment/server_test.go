package assignment_test

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-assignment-service/internal/config"
	"github.com/upassed/upassed-assignment-service/internal/handling"
	"github.com/upassed/upassed-assignment-service/internal/logging"
	"github.com/upassed/upassed-assignment-service/internal/middleware/common/auth"
	"github.com/upassed/upassed-assignment-service/internal/server"
	"github.com/upassed/upassed-assignment-service/internal/util"
	"github.com/upassed/upassed-assignment-service/internal/util/mocks"
	"github.com/upassed/upassed-assignment-service/pkg/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log"
	"os"
	"path/filepath"
	"testing"
)

var (
	assignmentClient client.AssignmentClient
	assignmentSvc    *mocks.AssignmentService
)

func TestMain(m *testing.M) {
	currentDir, _ := os.Getwd()
	projectRoot, err := util.GetProjectRoot(currentDir)
	if err != nil {
		log.Fatal("error to get project root folder: ", err)
	}

	if err := os.Setenv(config.EnvConfigPath, filepath.Join(projectRoot, "config", "test.yml")); err != nil {
		log.Fatal(err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("cfg load error: ", err)
	}

	logger := logging.New(cfg.Env)
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	authClient := mocks.NewAuthClientMW(ctrl)
	authClient.EXPECT().AuthenticationUnaryServerInterceptor().Return(emptyAuthMiddleware())

	assignmentSvc = mocks.NewAssignmentService(ctrl)
	assignmentServer := server.New(server.AppServerCreateParams{
		Config:            cfg,
		Log:               logger,
		AssignmentService: assignmentSvc,
		AuthClient:        authClient,
	})

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	cc, err := grpc.NewClient(fmt.Sprintf(":%s", cfg.GrpcServer.Port), opts...)
	if err != nil {
		log.Fatal("error creating client connection", err)
	}

	assignmentClient = client.NewAssignmentClient(cc)
	go func() {
		if err := assignmentServer.Run(); err != nil {
			os.Exit(1)
		}
	}()

	exitCode := m.Run()
	assignmentServer.GracefulStop()
	os.Exit(exitCode)
}

func TestFindAssignmentByFormID_InvalidRequest(t *testing.T) {
	request := &client.AssignmentFindByFormIDRequest{
		FormId: "invalid_uuid",
	}

	_, err := assignmentClient.FindByFormID(context.Background(), request)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, codes.InvalidArgument, convertedError.Code())
}

func TestFindAssignmentByFormID_ServiceError(t *testing.T) {
	request := &client.AssignmentFindByFormIDRequest{
		FormId: uuid.NewString(),
	}

	expectedServiceError := handling.New("some service error", codes.NotFound)
	assignmentSvc.EXPECT().
		FindByFormID(gomock.Any(), uuid.MustParse(request.GetFormId())).
		Return(nil, expectedServiceError)

	_, err := assignmentClient.FindByFormID(context.Background(), request)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedServiceError.Error(), convertedError.Message())
	assert.Equal(t, expectedServiceError.Code(), convertedError.Code())
}

func TestFindAssignmentByFormID_HappyPath(t *testing.T) {
	request := &client.AssignmentFindByFormIDRequest{
		FormId: uuid.NewString(),
	}

	expectedServiceResponse := util.RandomBusinessFormAssignment()
	assignmentSvc.EXPECT().
		FindByFormID(gomock.Any(), uuid.MustParse(request.GetFormId())).
		Return(expectedServiceResponse, nil)

	response, err := assignmentClient.FindByFormID(context.Background(), request)
	require.NoError(t, err)

	assert.Equal(t, len(expectedServiceResponse.GroupIDs), len(response.GetGroupIds()))
	for idx, groupID := range expectedServiceResponse.GroupIDs {
		assert.Equal(t, groupID.String(), response.GetGroupIds()[idx])
	}
}

func TestFindAssignmentByGroupID_InvalidRequest(t *testing.T) {
	request := &client.AssignmentFindByGroupIDRequest{
		GroupId: "invalid_uuid",
	}

	_, err := assignmentClient.FindByGroupID(context.Background(), request)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, codes.InvalidArgument, convertedError.Code())
}

func TestFindAssignmentByGroupID_ServiceError(t *testing.T) {
	request := &client.AssignmentFindByGroupIDRequest{
		GroupId: uuid.NewString(),
	}

	expectedServiceError := handling.New("some service error", codes.NotFound)
	assignmentSvc.EXPECT().
		FindByGroupID(gomock.Any(), uuid.MustParse(request.GetGroupId())).
		Return(nil, expectedServiceError)

	_, err := assignmentClient.FindByGroupID(context.Background(), request)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedServiceError.Error(), convertedError.Message())
	assert.Equal(t, expectedServiceError.Code(), convertedError.Code())
}

func TestFindAssignmentByGroupID_HappyPath(t *testing.T) {
	request := &client.AssignmentFindByGroupIDRequest{
		GroupId: uuid.NewString(),
	}

	expectedServiceResponse := util.RandomBusinessGroupAssignment()
	assignmentSvc.EXPECT().
		FindByGroupID(gomock.Any(), uuid.MustParse(request.GetGroupId())).
		Return(expectedServiceResponse, nil)

	response, err := assignmentClient.FindByGroupID(context.Background(), request)
	require.NoError(t, err)

	assert.Equal(t, len(expectedServiceResponse.FormIDs), len(response.GetFormIds()))
	for idx, formID := range expectedServiceResponse.FormIDs {
		assert.Equal(t, formID.String(), response.GetFormIds()[idx])
	}
}

func emptyAuthMiddleware() func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		ctxWithUsername := context.WithValue(ctx, auth.UsernameKey, "teacherUsername")

		return handler(ctxWithUsername, req)
	}
}
