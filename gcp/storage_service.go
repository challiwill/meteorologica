package gcp

import "google.golang.org/api/storage/v1"

type storageService struct {
	service *storage.Service
}

func (s *storageService) Buckets(projectID string) (*storage.Buckets, error) {
	buckets, err := s.service.Buckets.List(projectID).Do()
	if err != nil {
		return nil, err
	}
	return buckets, nil
}
