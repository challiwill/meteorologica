package csv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
)

type CSVReader interface {
	Read() ([]string, error)
	ReadAll() ([][]string, error)
}

type ReaderCleaner struct {
	Reader  CSVReader
	Cleaner *Cleaner
}

func NewReaderCleaner(body io.Reader, rowLen ...int) (*ReaderCleaner, error) {
	var max, min int
	switch {
	case len(rowLen) == 0:
		return nil, errors.New("Please provide an expected row length")
	case len(rowLen) == 1:
		max = rowLen[0]
	case len(rowLen) == 2:
		max = rowLen[0]
		min = rowLen[1]
	default:
		return nil, errors.New("Too many length arguments provided to NewReaderCleaner()")
	}

	cleaner, err := NewCleaner(max, min)
	if err != nil {
		return nil, err
	}
	csvReader := csv.NewReader(body)
	csvReader.FieldsPerRecord = -1
	csvReader.LazyQuotes = true
	return &ReaderCleaner{
		Reader:  csvReader,
		Cleaner: cleaner,
	}, nil
}

func (rc *ReaderCleaner) Read() ([]string, error) {
	for {
		report, err := rc.Reader.Read()
		if err != nil {
			return nil, err
		}
		if rc.Cleaner.IsFilledRow(report) && rc.Cleaner.IsRegularLengthRow(report) {
			return report, nil
		}
	}
}

func (rc *ReaderCleaner) ReadAll() ([][]string, error) {
	reports, err := rc.Reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("ReadAll Failed: %s", err.Error())
	}
	if len(reports) == 0 {
		return nil, NewEmptyReportError("")
	}

	reports = rc.Cleaner.RemoveEmptyRows(reports)
	if len(reports) == 0 {
		return nil, NewEmptyReportError("removing empty rows")
	}
	reports = rc.Cleaner.RemoveShortAndTruncateLongRows(reports)
	if len(reports) == 0 {
		return nil, NewEmptyReportError("removing short rows")
	}

	return reports, nil
}
