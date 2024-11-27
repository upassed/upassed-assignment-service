package assignment_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-assignment-service/internal/caching"
	"github.com/upassed/upassed-assignment-service/internal/caching/assignment"
	"github.com/upassed/upassed-assignment-service/internal/config"
	"github.com/upassed/upassed-assignment-service/internal/logging"
	"github.com/upassed/upassed-assignment-service/internal/testcontainer"
	"github.com/upassed/upassed-assignment-service/internal/util"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

var (
	redisClient *assignment.RedisClient
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
		log.Fatal("unable to parse config: ", err)
	}

	ctx := context.Background()
	logger := logging.New(cfg.Env)

	redisTestcontainer, err := testcontainer.NewRedisTestcontainer(ctx, cfg)
	if err != nil {
		log.Fatal("unable to run redis testcontainer: ", err)
	}

	port, err := redisTestcontainer.Start(ctx)
	if err != nil {
		log.Fatal("unable to get a postgres testcontainer real port: ", err)
	}

	cfg.Redis.Port = strconv.Itoa(port)
	redis, err := caching.OpenRedisConnection(cfg, logger)
	if err != nil {
		log.Fatal("unable to open connections to redis: ", err)
	}

	redisClient = assignment.New(redis, cfg, logger)
	exitCode := m.Run()
	if err := redisTestcontainer.Stop(ctx); err != nil {
		log.Println("unable to stop redis testcontainer: ", err)
	}

	os.Exit(exitCode)
}

func TestSaveByFormID_HappyPath(t *testing.T) {
	formID := uuid.New()
	ctx := context.Background()

	assignmentsToSave := util.RandomDomainAssignments()
	for _, assignmentToSave := range assignmentsToSave {
		assignmentToSave.FormID = formID
	}

	existingAssignments, err := redisClient.GetByFormID(ctx, formID)
	require.NoError(t, err)
	assert.Equal(t, 0, len(existingAssignments))

	err = redisClient.SaveByFormID(ctx, assignmentsToSave)
	require.NoError(t, err)

	existingAssignments, err = redisClient.GetByFormID(ctx, formID)
	require.NoError(t, err)
	assert.Equal(t, len(assignmentsToSave), len(existingAssignments))
}

func TestGetByFormID_AssignmentsNotFound(t *testing.T) {
	formID := uuid.New()
	ctx := context.Background()

	existingAssignments, err := redisClient.GetByFormID(ctx, formID)
	require.NoError(t, err)
	assert.Equal(t, 0, len(existingAssignments))
}

func TestGetByGroupID_AssignmentsNotFound(t *testing.T) {
	groupID := uuid.New()
	ctx := context.Background()

	existingAssignments, err := redisClient.GetByGroupID(ctx, groupID)
	require.NoError(t, err)
	assert.Equal(t, 0, len(existingAssignments))
}

func TestSaveByGroupID_HappyPath(t *testing.T) {
	groupID := uuid.New()
	ctx := context.Background()

	assignmentsToSave := util.RandomDomainAssignments()
	for _, assignmentToSave := range assignmentsToSave {
		assignmentToSave.GroupID = groupID
	}

	existingAssignments, err := redisClient.GetByGroupID(ctx, groupID)
	require.NoError(t, err)
	assert.Equal(t, 0, len(existingAssignments))

	for idx, assignmentToSave := range assignmentsToSave {
		err = redisClient.AddByGroupID(ctx, assignmentToSave)
		require.NoError(t, err)

		existingAssignments, err = redisClient.GetByGroupID(ctx, groupID)
		require.NoError(t, err)
		assert.Equal(t, idx+1, len(existingAssignments))
	}

	existingAssignments, err = redisClient.GetByGroupID(ctx, groupID)
	require.NoError(t, err)
	assert.Equal(t, len(assignmentsToSave), len(existingAssignments))
}
