package aws_test

import (
	"github.com/challiwill/meteorologica/aws"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Aws", func() {
	var client *aws.Client

	BeforeEach(func() {
		client = new(aws.Client)
	})

	Describe("Name", func() {
		It("returns the IAAS name", func() {
			Expect(client.Name()).To(Equal("AWS"))
		})
	})

	XDescribe("GetBillingData", func() {
		var (
			usage []byte
			err   error
		)

		JustBeforeEach(func() {
			usage, err = client.GetBillingData()
		})

		Context("When AWS returns a billing file", func() {
			BeforeEach(func() {

			})

			It("does not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

})
