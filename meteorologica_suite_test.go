package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMeteorologica(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Meteorologica Suite")
}
