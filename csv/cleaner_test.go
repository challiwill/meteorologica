package csv_test

import (
	. "github.com/challiwill/meteorologica/csv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cleaner", func() {
	var (
		cleaner *Cleaner
		err     error
	)

	Describe("NewCleaner", func() {
		Context("with a negative expected length", func() {
			BeforeEach(func() {
				cleaner, err = NewCleaner(-1)
			})

			It("errors", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("positive integer"))
			})
		})

		Context("with a zero expected length", func() {
			BeforeEach(func() {
				cleaner, err = NewCleaner(0)
			})

			It("errors", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("positive integer"))
			})
		})

		Context("with a positive expected length", func() {
			BeforeEach(func() {
				cleaner, err = NewCleaner(1)
			})

			It("works", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(cleaner).NotTo(BeNil())
			})
		})
	})

	Describe("RemoveEmptyRows", func() {
		Context("with valid csv", func() {
			var (
				cleaned  CSV
				original CSV
			)

			BeforeEach(func() {
				cleaner, err = NewCleaner(3)
				Expect(err).NotTo(HaveOccurred())
				original = CSV{
					[]string{"first header", " second hearder", "third header"},
					[]string{"one value", "two value", "three value"},
					[]string{"", "", ""},
					[]string{"more", "is", "better"},
					[]string{"", "", ""},
					[]string{"", "one thing", ""},
				}
			})

			JustBeforeEach(func() {
				cleaned = cleaner.RemoveEmptyRows(original)
			})

			It("returns a csv without empty rows", func() {
				expectedCSV := CSV{
					[]string{"first header", " second hearder", "third header"},
					[]string{"one value", "two value", "three value"},
					[]string{"more", "is", "better"},
					[]string{"", "one thing", ""},
				}
				Expect(cleaned).To(Equal(expectedCSV))
			})
		})

		Context("with irregular csv", func() {
			var (
				cleaned  CSV
				original CSV
			)

			BeforeEach(func() {
				cleaner, err = NewCleaner(3)
				Expect(err).NotTo(HaveOccurred())
				original = CSV{
					[]string{"first header", " second hearder", "third header"},
					[]string{"fewer", "values"},
					[]string{"", ""},
					[]string{"more", "really", "is", "better"},
					[]string{"", "", "", "", ""},
				}
			})

			JustBeforeEach(func() {
				cleaned = cleaner.RemoveEmptyRows(original)
			})

			It("returns a csv without empty rows", func() {
				expectedCSV := CSV{
					[]string{"first header", " second hearder", "third header"},
					[]string{"fewer", "values"},
					[]string{"more", "really", "is", "better"},
				}
				Expect(cleaned).To(Equal(expectedCSV))
			})

		})
	})

	Describe("RemoveIrregularLengthRows", func() {
		Context("with valid csv", func() {
			var (
				cleaned  CSV
				original CSV
			)

			BeforeEach(func() {
				original = CSV{
					[]string{"first header", " second hearder", "third header"},
					[]string{"fewer", "values"},
					[]string{"", ""},
					[]string{"three", "is", "key"},
					[]string{"more", "really", "is", "better"},
					[]string{"", "", ""},
					[]string{"", "", "", "", ""},
				}
			})

			JustBeforeEach(func() {
				cleaned = cleaner.RemoveIrregularLengthRows(original)
			})

			It("removes rows that don't have the right length", func() {
				expectedCSV := CSV{
					[]string{"first header", " second hearder", "third header"},
					[]string{"three", "is", "key"},
					[]string{"", "", ""},
				}
				Expect(cleaned).To(Equal(expectedCSV))
			})
		})
	})
})
