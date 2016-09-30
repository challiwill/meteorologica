package csv_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCsv(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Csv Suite")
}
