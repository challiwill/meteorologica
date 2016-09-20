package gcp

import (
	"net/http"
	"os"

	"google.golang.org/api/storage/v1"
)

type storageService struct {
	service *storage.Service
}

func (s *storageService) DailyUsage(bucketName string, objectName string) (*http.Response, error) {
	return s.service.Objects.Get(bucketName, objectName).Download()
}

func (s *storageService) Insert(bucketName string, object *storage.Object, file *os.File) (*storage.Object, error) {
	return s.service.Objects.Insert(bucketName, object).Media(file).Do()
}
