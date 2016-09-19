package azure_test

import (
	"github.com/challiwill/meteorologica/azure"

	. "github.com/onsi/ginkgo"
)

var _ = Describe("UsageReader", func() {
	var usageReader azure.UsageReader
	BeforeEach(func() {
		usageReader = azure.UsageReader{}
	})

})
