package csv_test

import (
	"errors"

	. "github.com/challiwill/meteorologica/csv"
	"github.com/challiwill/meteorologica/csv/csvfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

//go:generate counterfeiter io.Reader
//go:generate counterfeiter github.com/gocarina/gocsv.CSVReader

var _ = Describe("ReaderCleaner", func() {
	var (
		rc     *ReaderCleaner
		body   *csvfakes.FakeReader
		rowLen int
		err    error
	)

	Describe("NewReaderCleaner", func() {
		JustBeforeEach(func() {
			rc, err = NewReaderCleaner(body, rowLen)
		})

		Context("with negative row length", func() {
			BeforeEach(func() {
				body = new(csvfakes.FakeReader)
				rowLen = -1
			})

			It("errors", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("positive integer"))
			})
		})

		Context("with valid body and row length", func() {
			BeforeEach(func() {
				body = new(csvfakes.FakeReader)
				rowLen = 1
			})

			It("works", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(rc).NotTo(BeNil())
			})
		})
	})

	Describe("Read", func() {
		var (
			report []string
			reader *csvfakes.FakeCSVReader
			err    error
		)

		BeforeEach(func() {
			cleaner, err := NewCleaner(4)
			Expect(err).NotTo(HaveOccurred())
			reader = new(csvfakes.FakeCSVReader)
			rc = &ReaderCleaner{
				Reader:  reader,
				Cleaner: cleaner,
			}
		})

		JustBeforeEach(func() {
			report, err = rc.Read()
		})

		Context("when the read returns a filled row of regular length", func() {
			BeforeEach(func() {
				reader.ReadReturns([]string{"a", "comma", "separated", "value"}, nil)
			})

			It("works", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns the row", func() {
				Expect(report).To(Equal([]string{"a", "comma", "separated", "value"}))
			})
		})

		Context("when the read returns a filled row of irregular length", func() {
			Context("and there are following rows", func() {
				BeforeEach(func() {
					reader.ReadStub = func() ([]string, error) {
						if reader.ReadCallCount() == 0 {
							return []string{"comma", "separated", "value"}, nil
						}
						return []string{"a", "comma", "separated", "value"}, nil
					}
				})

				It("works", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("returns the row", func() {
					Expect(report).To(Equal([]string{"a", "comma", "separated", "value"}))
				})
			})

			Context("and there are no following rows", func() {
				BeforeEach(func() {
					reader.ReadStub = func() ([]string, error) {
						if reader.ReadCallCount() == 0 {
							return []string{"comma", "separated", "value"}, nil
						}
						return nil, errors.New("EOF")
					}
				})

				It("errors", func() {
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Context("when the read returns a row with no contents", func() {
			Context("and there are following rows", func() {
				BeforeEach(func() {
					reader.ReadStub = func() ([]string, error) {
						if reader.ReadCallCount() == 0 {
							return []string{"", "", ""}, nil
						}
						return []string{"a", "comma", "separated", "value"}, nil
					}
				})

				It("works", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("returns the row", func() {
					Expect(report).To(Equal([]string{"a", "comma", "separated", "value"}))
				})
			})

			Context("and there are no following rows", func() {
				BeforeEach(func() {
					reader.ReadStub = func() ([]string, error) {
						if reader.ReadCallCount() == 0 {
							return []string{"", "", "", ""}, nil
						}
						return nil, errors.New("EOF")
					}
				})

				It("errors", func() {
					Expect(err).To(HaveOccurred())
				})
			})

		})

		Context("when the read returns an error", func() {
			BeforeEach(func() {
				reader.ReadReturns(nil, errors.New("a read error"))
			})

			It("errors", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("a read error"))
			})
		})
	})

	Describe("ReadAll", func() {
		var (
			reports [][]string
			reader  *csvfakes.FakeCSVReader
			err     error
		)

		BeforeEach(func() {
			cleaner, err := NewCleaner(3)
			Expect(err).NotTo(HaveOccurred())
			reader = new(csvfakes.FakeCSVReader)
			rc = &ReaderCleaner{
				Reader:  reader,
				Cleaner: cleaner,
			}
		})

		JustBeforeEach(func() {
			reports, err = rc.ReadAll()
		})

		Context("when the read returns a slice of string slices", func() {
			BeforeEach(func() {
				reader.ReadAllReturns([][]string{
					[]string{"a", "first", "string"},
					[]string{"", "", ""},
					[]string{"second", "string"},
					[]string{"this", "one", "is", "longer"},
					[]string{"a", "fourth", "string"},
				}, nil)
			})

			It("works", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("removes the empty and short length rows, truncated the longer ones", func() {
				Expect(reports).To(Equal([][]string{
					[]string{"a", "first", "string"},
					[]string{"this", "one", "is"},
					[]string{"a", "fourth", "string"},
				}))
			})
		})

		Context("when all rows are short", func() {
			BeforeEach(func() {
				reader.ReadAllReturns([][]string{
					[]string{"a", "string"},
					[]string{"shorty"},
				}, nil)
			})

			It("returns a helpful error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Removing short rows resulted in empty report"))
			})
		})

		Context("when all rows are empty", func() {
			BeforeEach(func() {
				reader.ReadAllReturns([][]string{
					[]string{"", "", ""},
					[]string{"", "", ""},
				}, nil)
			})

			It("returns a helpful error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Removing empty rows resulted in empty report"))
			})
		})

		Context("when the read returns nil", func() {
			It("errors", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Report is empty"))
			})
		})

		Context("when the read returns an error", func() {
			BeforeEach(func() {
				reader.ReadAllReturns(nil, errors.New("a read error"))
			})

			It("errors", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("a read error"))
			})
		})
	})
})
