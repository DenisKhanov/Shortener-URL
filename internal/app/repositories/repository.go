package repositories

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DenisKhanov/shorterURL/internal/app/models"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

// URLInFileRepo auxiliary structure for serialization in jSON for save to file
type URLInFileRepo struct {
	UserID      uint32 `json:"user_id"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type URLInMemoryRepo struct {
	shortToOrigURL  map[string]string
	origToShortURL  map[string]string
	usersURLS       map[uint32][]models.URL
	batchBuffer     []URLInFileRepo
	batchCounter    uint8
	batchSize       uint8
	storageFilePath string
}

func NewURLInMemoryRepo(storageFilePath string) *URLInMemoryRepo {
	storage := URLInMemoryRepo{
		shortToOrigURL:  make(map[string]string),
		origToShortURL:  make(map[string]string),
		usersURLS:       make(map[uint32][]models.URL),
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
func (m *URLInMemoryRepo) StoreURLInDB(ctx context.Context, originalURL, shortURL string) error {
	userID, ok := ctx.Value("userID").(uint32)
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
func (m *URLInMemoryRepo) GetOriginalURLFromDB(ctx context.Context, shortURL string) (string, error) {
	originalURL, exists := m.shortToOrigURL[shortURL]
	if !exists {
		return "", errors.New("original URL not found")
	}
	return originalURL, nil
}
func (m *URLInMemoryRepo) GetShortURLFromDB(ctx context.Context, originalURL string) (string, error) {
	shortURL, exists := m.origToShortURL[originalURL]
	if !exists {
		return "", errors.New("short URL not found")
	}
	return shortURL, nil
}
func (m *URLInMemoryRepo) StoreBatchURLInDB(ctx context.Context, batchURLtoStores map[string]string) error {
	for shortURL, originalURL := range batchURLtoStores {
		if err := m.StoreURLInDB(ctx, originalURL, shortURL); err != nil {
			fmt.Println("Error saving to memory")
			return err
		}
	}
	return nil
}
func (m *URLInMemoryRepo) GetShortBatchURLFromDB(ctx context.Context, batchURLRequests []models.URLRequest) (map[string]string, error) {
	var shortsURL = make(map[string]string, len(batchURLRequests))

	for _, request := range batchURLRequests {
		if shortURL, ok := m.origToShortURL[request.OriginalURL]; ok {
			shortsURL[request.OriginalURL] = shortURL
		}
	}

	return shortsURL, nil
}
func (m *URLInMemoryRepo) GetUserURLSFromDB(ctx context.Context) ([]models.URL, error) {
	userID, ok := ctx.Value("userID").(uint32)
	if !ok {
		logrus.Errorf("context value is not userID: %v", userID)
	}
	userURLS, exists := m.usersURLS[userID]
	if !exists {
		return nil, errors.New("userID not found")
	}
	return userURLS, nil
}
