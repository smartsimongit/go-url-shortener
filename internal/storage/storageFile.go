package storage

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

type consumer struct {
	file    *os.File
	scanner *bufio.Scanner
}
type producer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewConsumer(fileName string) (*consumer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &consumer{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}
func (c *consumer) ReadShortURLs() (*URLRecords, error) {
	if !c.scanner.Scan() {
		return nil, c.scanner.Err()
	}
	data := c.scanner.Bytes()
	shortURLs := &URLRecords{}
	err := json.Unmarshal(data, &shortURLs)
	if err != nil {
		return nil, err
	}
	return shortURLs, nil
}
func (c *consumer) Close() error {
	return c.file.Close()
}

func NewProducer(fileName string) (*producer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return &producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}
func (p *producer) WriteShortURLs(shortURLs *URLRecords) error {
	return p.encoder.Encode(&shortURLs)
}
func (p *producer) Close() error {
	return p.file.Close()
}

func saveSortURLs(fileForSave string, shortURLs *URLRecords) bool {
	if fileForSave != "" {
		producer, err := NewProducer(fileForSave)
		if err != nil {
			log.Fatal(err)
		}
		defer producer.Close()
		if err := producer.WriteShortURLs(shortURLs); err != nil {
			log.Fatal(err)
		}
		return true
	}
	return false
}
func restoreShortURLs(fileForSave string) *URLRecords {
	consumer, err := NewConsumer(fileForSave)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()
	shortURLs, err := consumer.ReadShortURLs()
	if err != nil {
		log.Fatal(err)
	}
	return shortURLs
}
