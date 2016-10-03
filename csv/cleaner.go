package csv

import (
	"errors"
	"strings"
)

type Cleaner struct {
	expectedRowLength int
}

func NewCleaner(expectedLen int) (*Cleaner, error) {
	if expectedLen < 1 {
		return nil, errors.New("The expected row length must be a positive integer greater than zero")
	}
	return &Cleaner{expectedRowLength: expectedLen}, nil
}

func (c *Cleaner) RemoveEmptyRows(original CSV) CSV {
	cleaned := CSV{}
	for _, row := range original {
		if c.IsFilledRow(row) {
			cleaned = append(cleaned, row)
		}
	}

	return cleaned
}

func (c *Cleaner) RemoveIrregularLengthRows(original CSV) CSV {
	cleaned := CSV{}
	for _, row := range original {
		if c.IsRegularLengthRow(row) {
			cleaned = append(cleaned, row)
		}
	}

	return cleaned
}

func (c *Cleaner) IsRegularLengthRow(row []string) bool {
	return len(row) == c.expectedRowLength
}

func (c *Cleaner) IsFilledRow(row []string) bool {
	notEmpty := false
	for _, record := range row {
		if isNotEmptyString(record) {
			notEmpty = true
			break
		}
	}
	return notEmpty
}

func isNotEmptyString(test string) bool {
	trimmed := strings.TrimSpace(test)
	return trimmed != ""
}
