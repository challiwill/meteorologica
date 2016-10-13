package csv

import (
	"fmt"
	"strings"

	"github.com/challiwill/meteorologica/errare"
)

type Cleaner struct {
	maxRowLen int
	minRowLen int
}

func NewCleaner(maxRowLen int, minRowSlice ...int) (*Cleaner, error) {
	if maxRowLen < 1 {
		return nil, errare.NewCreationError("Cleaner", "The expected row length must be a positive integer greater than zero")
	}
	minRowLen := maxRowLen
	if len(minRowSlice) == 1 {
		minRowLen = minRowSlice[0]
		if minRowLen > maxRowLen {
			return nil, errare.NewCreationError("Cleaner", fmt.Sprintf("The minimum row length '%d' cannot be greater than the maximum row length '%d'", minRowLen, maxRowLen))
		}
	} else if len(minRowSlice) > 1 {
		return nil, errare.NewCreationError("Cleaner", "Too many arguments provided, requires only maximum row length and minimum row length")
	}
	return &Cleaner{
		maxRowLen: maxRowLen,
		minRowLen: minRowLen,
	}, nil
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

func (c *Cleaner) RemoveShortAndTruncateLongRows(original CSV) CSV {
	cleaned := CSV{}
	for _, row := range original {
		if len(row) >= c.maxRowLen {
			cleaned = append(cleaned, row[:c.maxRowLen])
		} else if len(row) >= c.minRowLen {
			cleaned = append(cleaned, row)
		}
	}

	return cleaned
}

func (c *Cleaner) TruncateRows(original CSV) CSV {
	cleaned := CSV{}
	for _, row := range original {
		if len(row) >= c.maxRowLen {
			cleaned = append(cleaned, row[:c.maxRowLen])
			continue
		}
		cleaned = append(cleaned, row)
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
	return len(row) <= c.maxRowLen && len(row) >= c.minRowLen
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
