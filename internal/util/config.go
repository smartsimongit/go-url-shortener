package util

import (
	"flag"
	"os"
)

var (
	serverAddress  = "SERVER_ADDRESS"
	baseURL        = "BASE_URL"
	fileStorageURL = "FILE_STORAGE_PATH"
)

type appConfigStruct struct {
	serverAddressValue  string
	baseURLValue        string
	fileStorageURLValue string
}

var appConfig appConfigStruct

func ConfigApp() {
	appConfig.serverAddressValue = os.Getenv(serverAddress)
	if appConfig.serverAddressValue == "" {
		flag.StringVar(&appConfig.serverAddressValue, "a", ":8080", "application start address")
	}

	appConfig.baseURLValue = os.Getenv(baseURL)
	if appConfig.baseURLValue == "" {
		flag.StringVar(&appConfig.baseURLValue, "b", "http://localhost:8080", "the base address of the resulting shortened URL")
	}

	appConfig.fileStorageURLValue = os.Getenv(fileStorageURL)
	if appConfig.fileStorageURLValue == "" {
		flag.StringVar(&appConfig.fileStorageURLValue, "f", "", "path to file with shorted URLs")
	}

}
