package repositoryes

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

type URLInFileRepo struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
type URLInMemoryRepo struct {
	shortToOrigURL  map[string]string
	origToShortURL  map[string]string
	lastUUID        int
	batchBuffer     []URLInFileRepo
	batchCounter    int
	batchSize       int
	storageFilePath string
}

func NewURLInMemoryRepo(storageFilePath string) *URLInMemoryRepo {
	storage := URLInMemoryRepo{
		shortToOrigURL:  make(map[string]string),
		origToShortURL:  make(map[string]string),
		lastUUID:        0,
		batchBuffer:     make([]URLInFileRepo, 0),
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
	}
	if err = scanner.Err(); err != nil {
		logrus.Error(err)
		return err
	}
	m.lastUUID, err = strconv.Atoi(bufferJSON.UUID)
	if err != nil {
		return err
	}
	return nil
}
func (m *URLInMemoryRepo) StoreURLInDB(originalURL, shortURL string) error {
	m.origToShortURL[originalURL] = shortURL
	m.shortToOrigURL[shortURL] = originalURL
	m.lastUUID++
	newUUID := strconv.Itoa(m.lastUUID)
	record := URLInFileRepo{
		UUID:        newUUID,
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
func (m *URLInMemoryRepo) GetOriginalURLFromDB(shortURL string) (string, error) {
	originalURL, exists := m.shortToOrigURL[shortURL]
	if !exists {
		return "", errors.New("original URL not found")
	}
	return originalURL, nil
}
func (m *URLInMemoryRepo) GetShortURLFromDB(originalURL string) (string, error) {
	shortURL, exists := m.origToShortURL[originalURL]
	if !exists {
		return "", errors.New("short URL not found")
	}
	return shortURL, nil
}
