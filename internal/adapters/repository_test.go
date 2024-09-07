package adapters

import (
	"context"
	"go.uber.org/zap"
	"net/url"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/sashaaro/url-shortener/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestMemURLRepository(t *testing.T) {
	repo := NewMemURLRepository()

	// Test Data
	testURL, _ := url.Parse("https://example.com")
	hashKey := domain.HashKey("short123")
	userID := uuid.New()

	// Add a URL
	err := repo.Add(context.Background(), hashKey, *testURL, userID)
	require.NoError(t, err, "should not return an error on Add")

	// Fetch the URL by hash key
	storedURL, err := repo.GetByHash(context.Background(), hashKey)
	require.NoError(t, err, "should not return an error on GetByHash")
	require.NotNil(t, storedURL, "stored URL should not be nil")
	require.Equal(t, testURL.String(), storedURL.String(), "stored URL should match the original")

	// Fetch URLs by user
	urlEntries, err := repo.GetByUser(context.Background(), userID)
	require.NoError(t, err, "should not return an error on GetByUser")
	require.Len(t, urlEntries, 1, "there should be one URL for this user")
	require.Equal(t, CreatePublicURL(hashKey), urlEntries[0].ShortURL, "short URL should match")
	require.Equal(t, testURL.String(), urlEntries[0].OriginalURL, "original URL should match")

	// Delete by user
	deleted, err := repo.DeleteByUser(context.Background(), []domain.HashKey{hashKey}, userID)
	require.NoError(t, err, "should not return an error on DeleteByUser")
	require.True(t, deleted, "should return true when URLs are deleted")

	// Ensure the URL is removed
	storedURL, err = repo.GetByHash(context.Background(), hashKey)
	require.NoError(t, err, "should not return an error on GetByHash after deletion")
	require.Nil(t, storedURL, "stored URL should be nil after deletion")
}

func TestFileURLRepository(t *testing.T) {
	// Setup temporary file for testing
	tempFile, err := os.CreateTemp("", "url_repo_test_*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name()) // Clean up after test

	// Create a logger for testing
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	logger, _ := config.Build()
	sugarLogger := logger.Sugar()

	// Setup wrapped in-memory repository and file repository
	memRepo := NewMemURLRepository()
	fileRepo := NewFileURLRepository(tempFile.Name(), memRepo, *sugarLogger)

	// Test Data
	testURL, _ := url.Parse("https://example.com")
	hashKey := domain.HashKey("short123")
	userID := uuid.New()

	// Add a URL
	err = fileRepo.Add(context.Background(), hashKey, *testURL, userID)
	require.NoError(t, err, "should not return an error on Add")

	// Fetch the URL by hash key
	storedURL, err := fileRepo.GetByHash(context.Background(), hashKey)
	require.NoError(t, err, "should not return an error on GetByHash")
	require.NotNil(t, storedURL, "stored URL should not be nil")
	require.Equal(t, testURL.String(), storedURL.String(), "stored URL should match the original")

	// Fetch URLs by user
	urlEntries, err := fileRepo.GetByUser(context.Background(), userID)
	require.NoError(t, err, "should not return an error on GetByUser")
	require.Len(t, urlEntries, 1, "there should be one URL for this user")
	require.Equal(t, CreatePublicURL(hashKey), urlEntries[0].ShortURL, "short URL should match")
	require.Equal(t, testURL.String(), urlEntries[0].OriginalURL, "original URL should match")

	// Delete by user
	deleted, err := fileRepo.DeleteByUser(context.Background(), []domain.HashKey{hashKey}, userID)
	require.NoError(t, err, "should not return an error on DeleteByUser")
	require.True(t, deleted, "should return true when URLs are deleted")

	// Ensure the URL is removed
	storedURL, err = fileRepo.GetByHash(context.Background(), hashKey)
	require.NoError(t, err, "should not return an error on GetByHash after deletion")
	require.Nil(t, storedURL, "stored URL should be nil after deletion")

	// Close the repository and file
	err = fileRepo.Close()
	require.NoError(t, err, "should not return an error on Close")

	// Reload the repository from the file to ensure persistence works
	reloadedRepo := NewFileURLRepository(tempFile.Name(), NewMemURLRepository(), *sugarLogger)
	urlEntries, err = reloadedRepo.GetByUser(context.Background(), userID)
	require.NoError(t, err, "should not return an error on GetByUser after reload")
	require.Len(t, urlEntries, 0, "there should be no URLs after reload since deletion occurred")
}
