package services

import (
	"flag"
	"os"
)

var (
	serverAddress  = "SERVER_ADDRESS"
	baseURL        = "BASE_URL"
	fileStorageURL = "FILE_STORAGE_PATH"
	bdAddressURL   = "DATABASE_DSN"
)

type appConfigStruct struct {
	ServerAddressValue  string
	BaseURLValue        string
	FileStorageURLValue string
	BdAddressURL        string
}

var AppConfig appConfigStruct

func ConfigApp() {

	AppConfig.BdAddressURL = os.Getenv(bdAddressURL)
	if AppConfig.BdAddressURL == "" {
		flag.StringVar(&AppConfig.BdAddressURL, "d", ":8080", "bd address url")
	}

	AppConfig.ServerAddressValue = os.Getenv(serverAddress)
	if AppConfig.ServerAddressValue == "" {
		flag.StringVar(&AppConfig.ServerAddressValue, "a", ":8080", "application start address")
	}

	AppConfig.BaseURLValue = os.Getenv(baseURL)
	if AppConfig.BaseURLValue == "" {
		flag.StringVar(&AppConfig.BaseURLValue, "b", "http://localhost:8080", "the base address of the resulting shortened URL")
	}

	AppConfig.FileStorageURLValue = os.Getenv(fileStorageURL)
	if AppConfig.FileStorageURLValue == "" {
		flag.StringVar(&AppConfig.FileStorageURLValue, "f", "storage.txt", "path to file with shorted URLs")
	}
	flag.Parse()
}
