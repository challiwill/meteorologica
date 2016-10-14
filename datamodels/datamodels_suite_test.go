package datamodels_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDatamodels(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Datamodels Suite")
}
