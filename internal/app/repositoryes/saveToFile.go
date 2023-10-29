package repositoryes

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type URLInFileRepo struct {
	UUID            string `json:"uuid"`
	ShortURL        string `json:"short_url"`
	OriginalURL     string `json:"original_url"`
	storageFilePath string
}

var countID int

func (r *URLInFileRepo) LoadLastUUID() error {
	file, err := os.Open(r.storageFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var lastLine []byte
	for scanner.Scan() {
		lastLine = scanner.Bytes()
	}
	if err = scanner.Err(); err != nil {
		return err
	}
	var lastRecord URLInFileRepo
	err = json.Unmarshal(lastLine, &lastRecord)
	if err != nil {
		return err
	}
	countID, err = strconv.Atoi(lastRecord.UUID)
	if err != nil {
		return err
	}
	return nil
}

func NewURLInFileRepo(storageFilePath string) *URLInFileRepo {
	dataFile := &URLInFileRepo{
		UUID:            "",
		ShortURL:        "",
		OriginalURL:     "",
		storageFilePath: storageFilePath,
	}
	dataFile.LoadLastUUID()
	logrus.Info(dataFile)
	return dataFile
}
func writeLine(file *os.File, data []byte) (n int, err error) {
	// Добавление переноса строки к срезу байт
	dataWithNewline := append(data, '\n')
	// Запись данных в файл
	return file.Write(dataWithNewline)
}

func (r *URLInFileRepo) StoreURLSInDB(originalURL, shortURL string) error {
	countID++
	r.UUID = strconv.Itoa(countID)
	r.ShortURL = shortURL
	r.OriginalURL = originalURL
	dir := filepath.Dir(r.storageFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatal(err)
	}
	file, err := os.OpenFile(r.storageFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logrus.Info(err)
		return err
	}
	defer file.Close()
	var data []byte
	data, err = json.Marshal(&r)
	if err != nil {
		return err
	}
	_, err = writeLine(file, data)
	if err != nil {
		return err
	}
	return nil
}
func (r *URLInFileRepo) GetOriginalURLFromDB(shortURL string) (string, error) {
	file, err := os.OpenFile(r.storageFilePath, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println(err)
		logrus.Info(err)
		return "", err

	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var data URLInFileRepo
		line := scanner.Bytes()
		err := json.Unmarshal(line, &data)
		if err != nil {
			fmt.Println(err)
			logrus.Info(err)
			return "", err // возвращает ошибку, если не может разобрать строку
		}
		fmt.Println(data)
		fmt.Println(data.ShortURL)
		if data.ShortURL == shortURL {

			return data.OriginalURL, nil
		}
	}
	return "", errors.New("URL not found")

}
func (r *URLInFileRepo) GetShortURLFromDB(originalURL string) (string, error) {
	file, err := os.Open(r.storageFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var data URLInFileRepo
		line := scanner.Bytes()
		err := json.Unmarshal(line, &data)
		if err != nil {
			return "", err // возвращает ошибку, если не может разобрать строку
		}
		fmt.Println(data)
		if data.OriginalURL == originalURL {
			return data.ShortURL, nil
		}
	}
	return "", errors.New("URL not found")

}
