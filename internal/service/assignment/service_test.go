package assignment_test

import (
	"context"
	"errors"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-assignment-service/internal/config"
	"github.com/upassed/upassed-assignment-service/internal/logging"
	"github.com/upassed/upassed-assignment-service/internal/middleware/common/auth"
	assignmentSvc "github.com/upassed/upassed-assignment-service/internal/service/assignment"
	"github.com/upassed/upassed-assignment-service/internal/util"
	"github.com/upassed/upassed-assignment-service/internal/util/mocks"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"os"
	"path/filepath"
	"testing"
)

var (
	cfg        *config.Config
	repository *mocks.AssignmentRepository
	service    assignmentSvc.Service
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

	cfg, err = config.Load()
	if err != nil {
		log.Fatal("unable to parse config: ", err)
	}

	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	repository = mocks.NewAssignmentRepository(ctrl)
	service = assignmentSvc.New(cfg, logging.New(config.EnvTesting), repository)

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestCreate_ErrorCheckingAssignmentDuplicates(t *testing.T) {
	teacherUsername := gofakeit.Username()
	ctx := context.WithValue(context.Background(), auth.UsernameKey, teacherUsername)

	assignmentToCreate := util.RandomBusinessFormAssignment()
	expectedRepositoryError := errors.New("some repo error")

	repository.EXPECT().
		CheckDuplicates(gomock.Any(), gomock.Any()).
		Return(nil, expectedRepositoryError)

	_, err := service.Create(ctx, assignmentToCreate)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepositoryError.Error(), convertedError.Message())
	assert.Equal(t, codes.Internal, convertedError.Code())
}

func TestCreate_ErrorDuplicateAssignmentsFound(t *testing.T) {
	teacherUsername := gofakeit.Username()
	ctx := context.WithValue(context.Background(), auth.UsernameKey, teacherUsername)

	assignmentToCreate := util.RandomBusinessFormAssignment()
	duplicateAssignments := util.RandomDomainAssignments()

	repository.EXPECT().
		CheckDuplicates(gomock.Any(), gomock.Any()).
		Return(duplicateAssignments, nil)

	_, err := service.Create(ctx, assignmentToCreate)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, assignmentSvc.ErrDuplicateAssignmentsFound.Error(), convertedError.Message())
	assert.Equal(t, codes.AlreadyExists, convertedError.Code())
}

func TestCreate_ErrorSavingAssignment(t *testing.T) {
	teacherUsername := gofakeit.Username()
	ctx := context.WithValue(context.Background(), auth.UsernameKey, teacherUsername)

	assignmentToCreate := util.RandomBusinessFormAssignment()
	expectedRepositoryError := errors.New("some repo error")

	repository.EXPECT().
		CheckDuplicates(gomock.Any(), gomock.Any()).
		Return(nil, nil)

	repository.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(expectedRepositoryError)

	_, err := service.Create(ctx, assignmentToCreate)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepositoryError.Error(), convertedError.Message())
	assert.Equal(t, codes.Internal, convertedError.Code())
}

func TestCreate_ErrorDeadlineExceeded(t *testing.T) {
	oldTimeout := cfg.Timeouts.EndpointExecutionTimeoutMS
	cfg.Timeouts.EndpointExecutionTimeoutMS = "0"

	teacherUsername := gofakeit.Username()
	ctx := context.WithValue(context.Background(), auth.UsernameKey, teacherUsername)

	assignmentToCreate := util.RandomBusinessFormAssignment()

	repository.EXPECT().
		CheckDuplicates(gomock.Any(), gomock.Any()).
		Return(nil, nil)

	repository.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(nil)

	_, err := service.Create(ctx, assignmentToCreate)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, assignmentSvc.ErrAssignmentCreateDeadlineExceeded.Error(), convertedError.Message())
	assert.Equal(t, codes.DeadlineExceeded, convertedError.Code())

	cfg.Timeouts.EndpointExecutionTimeoutMS = oldTimeout
}

func TestCreate_HappyPath(t *testing.T) {
	teacherUsername := gofakeit.Username()
	ctx := context.WithValue(context.Background(), auth.UsernameKey, teacherUsername)

	assignmentToCreate := util.RandomBusinessFormAssignment()

	repository.EXPECT().
		CheckDuplicates(gomock.Any(), gomock.Any()).
		Return(nil, nil)

	repository.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(nil)

	assignmentCreateResponse, err := service.Create(ctx, assignmentToCreate)
	require.NoError(t, err)

	assert.Equal(t, len(assignmentToCreate.GroupIDs), len(assignmentCreateResponse.CreatedAssignmentIDs))
}

func TestFindByFormID_ErrorInRepository(t *testing.T) {
	teacherUsername := gofakeit.Username()
	ctx := context.WithValue(context.Background(), auth.UsernameKey, teacherUsername)

	formID := uuid.New()
	expectedRepositoryError := errors.New("some repo error")

	repository.EXPECT().
		FindByFormID(gomock.Any(), formID).
		Return(nil, expectedRepositoryError)

	_, err := service.FindByFormID(ctx, formID)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepositoryError.Error(), convertedError.Message())
	assert.Equal(t, codes.Internal, convertedError.Code())
}

func TestFindByFormID_ErrorDeadlineExceeded(t *testing.T) {
	oldTimeout := cfg.Timeouts.EndpointExecutionTimeoutMS
	cfg.Timeouts.EndpointExecutionTimeoutMS = "0"

	teacherUsername := gofakeit.Username()
	ctx := context.WithValue(context.Background(), auth.UsernameKey, teacherUsername)

	formID := uuid.New()
	domainAssignments := util.RandomDomainAssignments()

	repository.EXPECT().
		FindByFormID(gomock.Any(), formID).
		Return(domainAssignments, nil)

	_, err := service.FindByFormID(ctx, formID)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, assignmentSvc.ErrAssignmentFindByFormIDDeadlineExceeded.Error(), convertedError.Message())
	assert.Equal(t, codes.DeadlineExceeded, convertedError.Code())

	cfg.Timeouts.EndpointExecutionTimeoutMS = oldTimeout
}

func TestFindByFormID_HappyPath(t *testing.T) {
	teacherUsername := gofakeit.Username()
	ctx := context.WithValue(context.Background(), auth.UsernameKey, teacherUsername)

	formID := uuid.New()
	domainAssignments := util.RandomDomainAssignments()

	repository.EXPECT().
		FindByFormID(gomock.Any(), formID).
		Return(domainAssignments, nil)

	foundAssignment, err := service.FindByFormID(ctx, formID)
	require.NoError(t, err)

	for idx, domainAssignment := range domainAssignments {
		assert.Equal(t, domainAssignment.GroupID, foundAssignment.GroupIDs[idx])
	}
}

func TestFindByGroupID_ErrorInRepository(t *testing.T) {
	teacherUsername := gofakeit.Username()
	ctx := context.WithValue(context.Background(), auth.UsernameKey, teacherUsername)

	groupID := uuid.New()
	expectedRepositoryError := errors.New("some repo error")

	repository.EXPECT().
		FindByGroupID(gomock.Any(), groupID).
		Return(nil, expectedRepositoryError)

	_, err := service.FindByGroupID(ctx, groupID)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepositoryError.Error(), convertedError.Message())
	assert.Equal(t, codes.Internal, convertedError.Code())
}

func TestFindByGroupID_ErrorDeadlineExceeded(t *testing.T) {
	oldTimeout := cfg.Timeouts.EndpointExecutionTimeoutMS
	cfg.Timeouts.EndpointExecutionTimeoutMS = "0"

	teacherUsername := gofakeit.Username()
	ctx := context.WithValue(context.Background(), auth.UsernameKey, teacherUsername)

	groupID := uuid.New()
	domainAssignments := util.RandomDomainAssignments()

	repository.EXPECT().
		FindByGroupID(gomock.Any(), groupID).
		Return(domainAssignments, nil)

	_, err := service.FindByGroupID(ctx, groupID)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, assignmentSvc.ErrAssignmentFindByGroupIDDeadlineExceeded.Error(), convertedError.Message())
	assert.Equal(t, codes.DeadlineExceeded, convertedError.Code())

	cfg.Timeouts.EndpointExecutionTimeoutMS = oldTimeout
}

func TestFindByGroupID_HappyPath(t *testing.T) {
	teacherUsername := gofakeit.Username()
	ctx := context.WithValue(context.Background(), auth.UsernameKey, teacherUsername)

	groupID := uuid.New()
	domainAssignments := util.RandomDomainAssignments()

	repository.EXPECT().
		FindByGroupID(gomock.Any(), groupID).
		Return(domainAssignments, nil)

	foundAssignment, err := service.FindByGroupID(ctx, groupID)
	require.NoError(t, err)

	for idx, domainAssignment := range domainAssignments {
		assert.Equal(t, domainAssignment.FormID, foundAssignment.FormIDs[idx])
	}
}
