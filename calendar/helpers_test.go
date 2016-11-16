package calendar_test

import (
	"time"

	. "github.com/challiwill/meteorologica/calendar"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Calendar", func() {
	Describe("PadMonth", func() {
		It("works", func() {
			expectedInputsOutputsPairs := map[time.Month]string{
				time.January:   "01",
				time.February:  "02",
				time.March:     "03",
				time.April:     "04",
				time.May:       "05",
				time.June:      "06",
				time.July:      "07",
				time.August:    "08",
				time.September: "09",
				time.October:   "10",
				time.November:  "11",
				time.December:  "12",
			}
			for input, expectedOutput := range expectedInputsOutputsPairs {
				Expect(PadMonth(input)).To(Equal(expectedOutput))
			}
		})
	})

	XDescribe("YesterdaysDate", func() {
		It("works", func() {})
	})
})
