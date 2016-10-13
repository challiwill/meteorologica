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

	BeforeEach(func() {
		cleaner, err = NewCleaner(3, 2)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("NewCleaner", func() {
		Context("with a negative maximum length", func() {
			BeforeEach(func() {
				cleaner, err = NewCleaner(-1)
			})

			It("errors", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("positive integer"))
			})
		})

		Context("with a zero maximum length", func() {
			BeforeEach(func() {
				cleaner, err = NewCleaner(0)
			})

			It("errors", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("positive integer"))
			})
		})

		Context("with a positive maximum length", func() {
			BeforeEach(func() {
				cleaner, err = NewCleaner(1)
			})

			It("works", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(cleaner).NotTo(BeNil())
			})
		})

		Context("with a positive minimum length", func() {
			Context("less than the maximum length", func() {
				BeforeEach(func() {
					cleaner, err = NewCleaner(2, 1)
				})

				It("works", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(cleaner).NotTo(BeNil())
				})
			})

			Context("equal to the maximum length", func() {
				BeforeEach(func() {
					cleaner, err = NewCleaner(1, 1)
				})

				It("works", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(cleaner).NotTo(BeNil())
				})
			})

			Context("greater than the maximum length", func() {
				BeforeEach(func() {
					cleaner, err = NewCleaner(1, 2)
				})

				It("errors", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("greater than"))
				})
			})
		})

		Context("with a negative minimum length", func() {
			BeforeEach(func() {
				cleaner, err = NewCleaner(1, -1)
			})

			It("works", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(cleaner).NotTo(BeNil())
			})
		})

		Context("with too many arguments", func() {
			BeforeEach(func() {
				cleaner, err = NewCleaner(1, 2, 3)
			})

			It("errors", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Too many arguments"))
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

	Describe("RemoveShortAndTruncateLongRows", func() {
		Context("with valid csv", func() {
			var (
				cleaned  CSV
				original CSV
			)

			BeforeEach(func() {
				original = CSV{
					[]string{"three", "is", "key"},
					[]string{"fewer", "values"},
					[]string{"fewest"},
					[]string{"more", "really", "is", "better"},
					[]string{"", "", "", "", ""},
				}
			})

			JustBeforeEach(func() {
				cleaned = cleaner.RemoveShortAndTruncateLongRows(original)
			})

			It("removes rows that are too short and truncates rows that are too long", func() {
				expectedCSV := CSV{
					[]string{"three", "is", "key"},
					[]string{"fewer", "values"},
					[]string{"more", "really", "is"},
					[]string{"", "", ""},
				}
				Expect(cleaned).To(Equal(expectedCSV))
			})
		})
	})

	Describe("TruncateRows", func() {
		Context("with valid csv", func() {
			var (
				cleaned  CSV
				original CSV
			)

			BeforeEach(func() {
				original = CSV{
					[]string{"three", "is", "key"},
					[]string{"fewer", "values"},
					[]string{"more", "really", "is", "better"},
					[]string{"", "", "", "", ""},
				}
			})

			JustBeforeEach(func() {
				cleaned = cleaner.TruncateRows(original)
			})

			It("removes rows that don't have the right length", func() {
				expectedCSV := CSV{
					[]string{"three", "is", "key"},
					[]string{"fewer", "values"},
					[]string{"more", "really", "is"},
					[]string{"", "", ""},
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
					[]string{"fewest"},
					[]string{"", ""},
					[]string{""},
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
					[]string{"fewer", "values"},
					[]string{"", ""},
					[]string{"three", "is", "key"},
					[]string{"", "", ""},
				}
				Expect(cleaned).To(Equal(expectedCSV))
			})
		})
	})
})
