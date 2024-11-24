package assignment_test

import (
	"context"
	"errors"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/golang/mock/gomock"
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

	assignmentToCreate := util.RandomBusinessAssignment()
	expectedRepositoryError := errors.New("some repo error")

	repository.EXPECT().
		CheckDuplicates(gomock.Any(), gomock.Any()).
		Return(expectedRepositoryError)

	_, err := service.Create(ctx, assignmentToCreate)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepositoryError.Error(), convertedError.Message())
	assert.Equal(t, codes.Internal, convertedError.Code())
}

func TestCreate_ErrorSavingAssignment(t *testing.T) {
	teacherUsername := gofakeit.Username()
	ctx := context.WithValue(context.Background(), auth.UsernameKey, teacherUsername)

	assignmentToCreate := util.RandomBusinessAssignment()
	expectedRepositoryError := errors.New("some repo error")

	repository.EXPECT().
		CheckDuplicates(gomock.Any(), gomock.Any()).
		Return(nil)

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

	assignmentToCreate := util.RandomBusinessAssignment()

	repository.EXPECT().
		CheckDuplicates(gomock.Any(), gomock.Any()).
		Return(nil)

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

	assignmentToCreate := util.RandomBusinessAssignment()

	repository.EXPECT().
		CheckDuplicates(gomock.Any(), gomock.Any()).
		Return(nil)

	repository.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(nil)

	assignmentCreateResponse, err := service.Create(ctx, assignmentToCreate)
	require.NoError(t, err)

	assert.NotNil(t, assignmentCreateResponse.CreatedAssignmentID)
}