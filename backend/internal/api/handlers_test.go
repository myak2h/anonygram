package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"anonygram/internal/config"
	"anonygram/internal/models"
)

// Mocks

// ImageRepository Mock

type MockImageRepository struct {
	mock.Mock
}

func (m *MockImageRepository) Add(image models.Image) error {
	args := m.Called(image)
	return args.Error(0)
}

func (m *MockImageRepository) List(tags []string) []models.Image {
	args := m.Called(tags)
	return args.Get(0).([]models.Image)
}

// FileRepository Mock

type MockFileRepository struct {
	mock.Mock
}

func (m *MockFileRepository) Save(src io.Reader) (string, error) {
	args := m.Called(src)
	return args.String(0), args.Error(1)
}

// Broadcaster Mock

type MockBroadcaster struct {
	mock.Mock
}

func (m *MockBroadcaster) Broadcast(img models.Image) bool {
	args := m.Called(img)
	return args.Bool(0)
}

func (m *MockBroadcaster) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

// Test Helpers

func testConfig() *config.Config {
	return &config.Config{
		Port:           "8080",
		UploadPath:     "./uploads",
		AllowedOrigins: []string{"*"},
		MaxUploadSize:  10 << 20,
	}
}

func createMultipartRequest(fields map[string]string, fileName string, fileContent []byte) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, val := range fields {
		_ = writer.WriteField(key, val)
	}

	if fileName != "" {
		part, err := writer.CreateFormFile("image", fileName)
		if err != nil {
			return nil, err
		}
		_, _ = part.Write(fileContent)
	}

	_ = writer.Close()

	req, err := http.NewRequest("POST", "/upload", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

// Tests

func TestNewServer(t *testing.T) {
	imageRepo := &MockImageRepository{}
	fileRepo := &MockFileRepository{}
	cfg := testConfig()
	hub := &MockBroadcaster{}

	server := NewServer(imageRepo, fileRepo, cfg, hub)

	assert.NotNil(t, server)
	assert.Equal(t, imageRepo, server.imageRepo)
	assert.Equal(t, fileRepo, server.fileRepo)
	assert.Equal(t, cfg, server.config)
	assert.Equal(t, hub, server.hub)
}

func TestServer_ListImages(t *testing.T) {
	t.Run("no images", func(t *testing.T) {
		imageRepo := &MockImageRepository{}
		imageRepo.On("List", mock.Anything).Return([]models.Image{})
		fileRepo := &MockFileRepository{}
		hub := &MockBroadcaster{}

		server := NewServer(imageRepo, fileRepo, testConfig(), hub)
		req, _ := http.NewRequest("GET", "/images", nil)
		w := httptest.NewRecorder()

		server.Routes().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		imageRepo.AssertCalled(t, "List", mock.Anything)
		assert.JSONEq(t, "[]", w.Body.String())

	})

	t.Run("images exist", func(t *testing.T) {
		images := []models.Image{
			{ID: "1"},
			{ID: "2"},
		}
		imageRepo := &MockImageRepository{}
		imageRepo.On("List", []string(nil)).Return(images)
		fileRepo := &MockFileRepository{}
		hub := &MockBroadcaster{}

		server := NewServer(imageRepo, fileRepo, testConfig(), hub)
		req, _ := http.NewRequest("GET", "/images", nil)
		w := httptest.NewRecorder()

		server.Routes().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		imageRepo.AssertCalled(t, "List", []string(nil))

		responseImages := []models.Image{}
		err := json.Unmarshal(w.Body.Bytes(), &responseImages)
		assert.NoError(t, err)
		assert.Len(t, responseImages, 2)
	})

	t.Run("with tags filter", func(t *testing.T) {
		images := []models.Image{
			{ID: "1", Tags: []string{"nature"}},
		}
		imageRepo := &MockImageRepository{}
		imageRepo.On("List", []string{"nature", "sunset"}).Return(images)
		fileRepo := &MockFileRepository{}
		hub := &MockBroadcaster{}

		server := NewServer(imageRepo, fileRepo, testConfig(), hub)
		req, _ := http.NewRequest("GET", "/images?tag=nature&tag=sunset", nil)
		w := httptest.NewRecorder()

		server.Routes().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		imageRepo.AssertCalled(t, "List", []string{"nature", "sunset"})
	})
}

func TestServer_UploadImage(t *testing.T) {
	t.Run("missing title", func(t *testing.T) {
		imageRepo := &MockImageRepository{}
		fileRepo := &MockFileRepository{}
		hub := &MockBroadcaster{}

		server := NewServer(imageRepo, fileRepo, testConfig(), hub)
		req, _ := createMultipartRequest(map[string]string{}, "test.jpg", []byte("fake image"))
		w := httptest.NewRecorder()

		server.Routes().ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), ErrTitleRequired.Error())
	})

	t.Run("missing image file", func(t *testing.T) {
		imageRepo := &MockImageRepository{}
		fileRepo := &MockFileRepository{}
		hub := &MockBroadcaster{}

		server := NewServer(imageRepo, fileRepo, testConfig(), hub)
		req, _ := createMultipartRequest(map[string]string{"title": "Test Image"}, "", nil)
		w := httptest.NewRecorder()

		server.Routes().ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), ErrImageFileRequired.Error())
	})

	t.Run("successful upload", func(t *testing.T) {
		imageRepo := &MockImageRepository{}
		imageRepo.On("Add", mock.AnythingOfType("models.Image")).Return(nil)
		fileRepo := &MockFileRepository{}
		fileRepo.On("Save", mock.Anything).Return("/uploads/abc123.jpg", nil)
		hub := &MockBroadcaster{}
		hub.On("Broadcast", mock.AnythingOfType("models.Image")).Return(true)

		server := NewServer(imageRepo, fileRepo, testConfig(), hub)
		req, _ := createMultipartRequest(map[string]string{"title": "My Image", "tags": "nature, sunset"}, "test.jpg", []byte("fake image"))
		w := httptest.NewRecorder()

		server.Routes().ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		fileRepo.AssertCalled(t, "Save", mock.Anything)
		imageRepo.AssertCalled(t, "Add", mock.AnythingOfType("models.Image"))
		hub.AssertCalled(t, "Broadcast", mock.AnythingOfType("models.Image"))

		var response models.Image
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "My Image", response.Title)
		assert.Equal(t, []string{"nature", "sunset"}, response.Tags)
		assert.Equal(t, "/uploads/abc123.jpg", response.URL)
		assert.NotEmpty(t, response.ID)
	})

	t.Run("file save error", func(t *testing.T) {
		imageRepo := &MockImageRepository{}
		fileRepo := &MockFileRepository{}
		fileRepo.On("Save", mock.Anything).Return("", errors.New("disk full"))
		hub := &MockBroadcaster{}

		server := NewServer(imageRepo, fileRepo, testConfig(), hub)
		req, _ := createMultipartRequest(map[string]string{"title": "My Image"}, "test.jpg", []byte("fake image"))
		w := httptest.NewRecorder()

		server.Routes().ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), ErrFileSaveFailed.Error())
	})

	t.Run("image metadata save error", func(t *testing.T) {
		imageRepo := &MockImageRepository{}
		imageRepo.On("Add", mock.AnythingOfType("models.Image")).Return(errors.New("db error"))
		fileRepo := &MockFileRepository{}
		fileRepo.On("Save", mock.Anything).Return("/uploads/abc123.jpg", nil)
		hub := &MockBroadcaster{}

		server := NewServer(imageRepo, fileRepo, testConfig(), hub)
		req, _ := createMultipartRequest(map[string]string{"title": "My Image"}, "test.jpg", []byte("fake image"))
		w := httptest.NewRecorder()

		server.Routes().ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), ErrMetadataSaveFailed.Error())
	})
}
