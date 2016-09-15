package gcp

import (
	"net/http"

	"google.golang.org/api/storage/v1"
)

type storageService struct {
	service *storage.Service
}

func (s *storageService) DailyUsage(bucketName string, objectName string) (*http.Response, error) {
	return s.service.Objects.Get(bucketName, objectName).Download()
}
