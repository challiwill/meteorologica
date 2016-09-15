package gcp

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	storage "google.golang.org/api/storage/v1"
)

//go:generate counterfeiter . StorageService

type StorageService interface {
	Buckets(string) (*storage.Buckets, error)
}

type DetailedUsageReport struct{}

type Client struct {
	StorageService StorageService
}

func NewClient(jsonCredentials []byte) (*Client, error) {
	jwtConfig, err := google.JWTConfigFromJSON(jsonCredentials, "https://www.googleapis.com/auth/devstorage.read_only")
	if err != nil {
		return nil, err
	}
	service, err := storage.New(jwtConfig.Client(oauth2.NoContext))
	if err != nil {
		return nil, err
	}
	return &Client{
		StorageService: &storageService{service: service},
	}, nil
}

func (c Client) MonthlyUsageReport() (DetailedUsageReport, error) {
	_, err := c.StorageService.Buckets("cma-test")
	if err != nil {
		return DetailedUsageReport{}, err
	}

	return DetailedUsageReport{}, nil
}
