package csv

import "strings"

type CSV [][]string

type Cleaner struct{}

func NewCleaner() *Cleaner {
	return new(Cleaner)
}

func (c *Cleaner) RemoveEmptyRows(original CSV) CSV {
	cleaned := CSV{}
	for _, row := range original {
		notEmpty := false
		for _, record := range row {
			if isNotEmptyString(record) {
				notEmpty = true
				break
			}
		}
		if notEmpty {
			cleaned = append(cleaned, row)
		}
	}

	return cleaned
}

func (c *Cleaner) RemoveIrregularLengthRows(original CSV, expectedLen int) CSV {
	cleaned := CSV{}
	for _, row := range original {
		if len(row) == expectedLen {
			cleaned = append(cleaned, row)
		}
	}

	return cleaned
}

func isNotEmptyString(test string) bool {
	trimmed := strings.TrimSpace(test)
	return trimmed != ""
}
