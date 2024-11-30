package assignment_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-assignment-service/internal/caching"
	cache "github.com/upassed/upassed-assignment-service/internal/caching/assignment"
	"github.com/upassed/upassed-assignment-service/internal/config"
	"github.com/upassed/upassed-assignment-service/internal/logging"
	"github.com/upassed/upassed-assignment-service/internal/repository"
	"github.com/upassed/upassed-assignment-service/internal/repository/assignment"
	"github.com/upassed/upassed-assignment-service/internal/testcontainer"
	"github.com/upassed/upassed-assignment-service/internal/util"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

var (
	assignmentCache      *cache.RedisClient
	assignmentRepository assignment.Repository
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
	postgresTestcontainer, err := testcontainer.NewPostgresTestcontainer(ctx)
	if err != nil {
		log.Fatal("unable to create a testcontainer: ", err)
	}

	port, err := postgresTestcontainer.Start(ctx)
	if err != nil {
		log.Fatal("unable to get a postgres testcontainer real port: ", err)
	}

	cfg.Storage.Port = strconv.Itoa(port)
	logger := logging.New(cfg.Env)
	if err := postgresTestcontainer.Migrate(cfg, logger); err != nil {
		log.Fatal("unable to run migrations: ", err)
	}

	redisTestcontainer, err := testcontainer.NewRedisTestcontainer(ctx, cfg)
	if err != nil {
		log.Fatal("unable to run redis testcontainer: ", err)
	}

	port, err = redisTestcontainer.Start(ctx)
	if err != nil {
		log.Fatal("unable to get a postgres testcontainer real port: ", err)
	}

	cfg.Redis.Port = strconv.Itoa(port)
	db, err := repository.OpenGormDbConnection(cfg, logger)
	if err != nil {
		log.Fatal("unable to open connection to postgres: ", err)
	}

	redis, err := caching.OpenRedisConnection(cfg, logger)
	if err != nil {
		log.Fatal("unable to open connection to redis: ", err)
	}

	assignmentCache = cache.New(redis, cfg, logger)
	assignmentRepository = assignment.New(db, redis, cfg, logger)
	exitCode := m.Run()
	if err := postgresTestcontainer.Stop(ctx); err != nil {
		log.Println("unable to stop postgres testcontainer: ", err)
	}

	if err := redisTestcontainer.Stop(ctx); err != nil {
		log.Println("unable to stop redis testcontainer: ", err)
	}

	os.Exit(exitCode)
}

func TestCheckDuplicates_DuplicatesNotFound(t *testing.T) {
	domainAssignments := util.RandomDomainAssignments()

	foundDuplicates, err := assignmentRepository.CheckDuplicates(context.Background(), domainAssignments)
	require.NoError(t, err)

	assert.Equal(t, 0, len(foundDuplicates))
}

func TestSave_HappyPath(t *testing.T) {
	ctx := context.Background()
	formID := uuid.New()
	assignmentsToSave := util.RandomDomainAssignments()

	for _, domainAssignment := range assignmentsToSave {
		domainAssignment.FormID = formID
	}

	err := assignmentRepository.Save(ctx, assignmentsToSave)
	require.NoError(t, err)

	for _, savedAssignment := range assignmentsToSave {
		cachedAssignments, err := assignmentCache.GetByGroupID(ctx, savedAssignment.GroupID)
		require.NoError(t, err)

		assert.Equal(t, 1, len(cachedAssignments))
		assert.Equal(t, formID, cachedAssignments[0].FormID)
	}

	cachedAssignments, err := assignmentCache.GetByFormID(ctx, formID)
	require.NoError(t, err)

	assert.Equal(t, len(assignmentsToSave), len(cachedAssignments))
}

func TestFindByFormID_HappyPath(t *testing.T) {
	ctx := context.Background()
	formID := uuid.New()
	assignmentsToSave := util.RandomDomainAssignments()

	for _, domainAssignment := range assignmentsToSave {
		domainAssignment.FormID = formID
	}

	err := assignmentRepository.Save(ctx, assignmentsToSave)
	require.NoError(t, err)

	foundAssignments, err := assignmentRepository.FindByFormID(ctx, formID)
	require.NoError(t, err)

	assert.Equal(t, len(assignmentsToSave), len(foundAssignments))
	for idx, foundAssignment := range foundAssignments {
		assert.Equal(t, assignmentsToSave[idx].ID, foundAssignment.ID)
		assert.Equal(t, assignmentsToSave[idx].FormID, foundAssignment.FormID)
		assert.Equal(t, assignmentsToSave[idx].GroupID, foundAssignment.GroupID)
	}

	cachedAssignments, err := assignmentCache.GetByFormID(ctx, formID)
	require.NoError(t, err)

	assert.Equal(t, len(assignmentsToSave), len(cachedAssignments))
	for idx, cachedAssignment := range cachedAssignments {
		assert.Equal(t, assignmentsToSave[idx].ID, cachedAssignment.ID)
		assert.Equal(t, assignmentsToSave[idx].FormID, cachedAssignment.FormID)
		assert.Equal(t, assignmentsToSave[idx].GroupID, cachedAssignment.GroupID)
	}
}

func TestFindByGroupID_HappyPath(t *testing.T) {
	ctx := context.Background()
	formID := uuid.New()
	groupID := uuid.New()
	assignmentsToSave := util.RandomDomainAssignments()

	for _, domainAssignment := range assignmentsToSave {
		domainAssignment.FormID = formID
	}

	assignmentsToSave[0].GroupID = groupID
	err := assignmentRepository.Save(ctx, assignmentsToSave)
	require.NoError(t, err)

	groupAssignments, err := assignmentRepository.FindByGroupID(ctx, groupID)
	require.NoError(t, err)

	assert.Equal(t, 1, len(groupAssignments))
	assert.Equal(t, assignmentsToSave[0], groupAssignments[0])

	cachedAssignments, err := assignmentCache.GetByGroupID(ctx, groupID)
	require.NoError(t, err)

	assert.Equal(t, 1, len(cachedAssignments))
	assert.Equal(t, assignmentsToSave[0], cachedAssignments[0])
}
