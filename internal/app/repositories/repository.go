// Package repositories provides implementations of data storage for managing shortened URLs.
// It includes functionality to store, retrieve, and delete URLs using in-memory and file-based storage.
package repositories

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DenisKhanov/shorterURL/internal/app/models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

// URLInFileRepo auxiliary structure for serialization in jSON for save to file
type URLInFileRepo struct {
	UserID      uuid.UUID `json:"user_id"`
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
}

// URLInMemoryRepo represents an in-memory repository for managing shortened URLs.
type URLInMemoryRepo struct {
	shortToOrigURL  map[string]string
	origToShortURL  map[string]string
	usersURLS       map[uuid.UUID][]models.URL
	batchBuffer     []URLInFileRepo
	batchCounter    uint8
	batchSize       uint8
	storageFilePath string
}

// NewURLInMemoryRepo creates a new instance of URLInMemoryRepo.
// It takes a file path for storing data.
func NewURLInMemoryRepo(storageFilePath string) *URLInMemoryRepo {
	storage := URLInMemoryRepo{
		shortToOrigURL:  make(map[string]string),
		origToShortURL:  make(map[string]string),
		usersURLS:       make(map[uuid.UUID][]models.URL),
		batchBuffer:     []URLInFileRepo{},
		batchCounter:    0,
		batchSize:       100,
		storageFilePath: storageFilePath,
	}
	err := storage.readFileToMemoryURL()
	if err != nil {
		logrus.Error(err)
	}
	return &storage
}

// readFileToMemoryURL read data from file and write it to memory (to URLInMemoryRepo)
func (m *URLInMemoryRepo) readFileToMemoryURL() error {
	file, err := os.Open(m.storageFilePath)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var buffer []byte
	var bufferJSON URLInFileRepo
	for scanner.Scan() {
		buffer = scanner.Bytes()
		err = json.Unmarshal(buffer, &bufferJSON)
		if err != nil {
			logrus.Error(err)
			return err
		}
		m.shortToOrigURL[bufferJSON.ShortURL] = bufferJSON.OriginalURL
		m.origToShortURL[bufferJSON.OriginalURL] = bufferJSON.ShortURL
		m.usersURLS[bufferJSON.UserID] = append(m.usersURLS[bufferJSON.UserID], models.URL{ShortURL: bufferJSON.ShortURL, OriginalURL: bufferJSON.OriginalURL})
	}
	if err = scanner.Err(); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

// SaveBatchToFile read data from the memory (URLInMemoryRepo batchBuffer) and write to the file in a batch operation
func (m *URLInMemoryRepo) SaveBatchToFile() error {
	startTime := time.Now() // Засекаем время начала операции
	file, err := os.OpenFile(m.storageFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	encoder := json.NewEncoder(writer)
	for _, v := range m.batchBuffer {
		err = encoder.Encode(v)
		if err != nil {
			return err
		}
	}
	err = writer.Flush() // Запись оставшихся данных из буфера в файл
	if err != nil {
		return err
	}

	elapsedTime := time.Since(startTime) // Вычисляем затраченное время
	logrus.Infof("%d URL saved in %v", m.batchCounter, elapsedTime)
	m.batchBuffer = make([]URLInFileRepo, 0, m.batchSize)
	return nil
}

// StoreURLInDB saves a mapping between an original URL and its shortened version in the database.
// It returns an error if the saving process fails.
func (m *URLInMemoryRepo) StoreURLInDB(ctx context.Context, originalURL, shortURL string) error {
	userID, ok := ctx.Value(models.UserIDKey).(uuid.UUID)
	if !ok {
		logrus.Errorf("context value is not userID: %v", userID)
	}
	m.origToShortURL[originalURL] = shortURL
	m.shortToOrigURL[shortURL] = originalURL
	m.usersURLS[userID] = append(m.usersURLS[userID], models.URL{ShortURL: shortURL, OriginalURL: originalURL})

	record := URLInFileRepo{
		UserID:      userID,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}
	m.batchBuffer = append(m.batchBuffer, record)
	m.batchCounter++
	if m.batchCounter >= m.batchSize {
		err := m.SaveBatchToFile()
		if err != nil {
			return err
		}
		logrus.Infof("%d URL saved to file", m.batchCounter)
		m.batchCounter = 0
	}

	if m.origToShortURL[originalURL] == "" || m.shortToOrigURL[shortURL] == "" {
		err := errors.New("error saving shortUrl")
		logrus.Error(err)
		return err
	}
	return nil
}

// GetOriginalURLFromDB retrieves the original URL corresponding to a given shortened URL from the database.
// It returns the original URL and any error encountered during the retrieval.
func (m *URLInMemoryRepo) GetOriginalURLFromDB(ctx context.Context, shortURL string) (string, error) {
	originalURL, exists := m.shortToOrigURL[shortURL]
	if !exists {
		return "", errors.New("original URL not found")
	}
	return originalURL, nil
}

// GetShortURLFromDB retrieves the shortened version of a given original URL from the database.
// It returns the shortened URL and any error encountered during the retrieval.
func (m *URLInMemoryRepo) GetShortURLFromDB(ctx context.Context, originalURL string) (string, error) {
	shortURL, exists := m.origToShortURL[originalURL]
	if !exists {
		return "", errors.New("short URL not found")
	}
	return shortURL, nil
}

// StoreBatchURLInDB saves multiple URL mappings in the database in a batch operation.
// The input is a map where keys are shortened URLs and values are the corresponding original URLs.
// It returns an error if the batch saving process fails.
func (m *URLInMemoryRepo) StoreBatchURLInDB(ctx context.Context, batchURLtoStores map[string]string) error {
	for shortURL, originalURL := range batchURLtoStores {
		if err := m.StoreURLInDB(ctx, originalURL, shortURL); err != nil {
			fmt.Println("Error saving to memory")
			return err
		}
	}
	return nil
}

// GetShortBatchURLFromDB retrieves multiple shortened URLs corresponding to a batch of original URLs from the database.
// The input is a slice of URLRequest objects containing original URLs.
//
//	It returns found in database a map of original URLs to their shortened counterparts and any error encountered during the retrieval.
func (m *URLInMemoryRepo) GetShortBatchURLFromDB(ctx context.Context, batchURLRequests []models.URLRequest) (map[string]string, error) {
	var shortsURL = make(map[string]string, len(batchURLRequests))

	for _, request := range batchURLRequests {
		if shortURL, ok := m.origToShortURL[request.OriginalURL]; ok {
			shortsURL[request.OriginalURL] = shortURL
		}
	}

	return shortsURL, nil
}

// GetUserURLSFromDB takes a slice of models.URL objects for a specific user from DB
func (m *URLInMemoryRepo) GetUserURLSFromDB(ctx context.Context) ([]models.URL, error) {
	userID, ok := ctx.Value(models.UserIDKey).(uuid.UUID)
	if !ok {
		logrus.Errorf("context value is not userID: %v", userID)
	}
	userURLS, exists := m.usersURLS[userID]
	if !exists {
		return nil, errors.New("userID not found")
	}
	return userURLS, nil
}

// MarkURLsAsDeleted marks user URLs as deleted in DB
func (m *URLInMemoryRepo) MarkURLsAsDeleted(ctx context.Context, URLSToDel []string) error {
	return nil
}

func (m *URLInMemoryRepo) Stats(ctx context.Context) (models.Stats, error) {
	urls := len(m.shortToOrigURL)
	users := len(m.usersURLS)
	stats := models.Stats{Urls: urls, Users: users}
	return stats, nil
}
