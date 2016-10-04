package csv

import (
	"encoding/csv"
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

func NewReaderCleaner(body io.Reader, rowLen int) (*ReaderCleaner, error) {
	cleaner, err := NewCleaner(rowLen)
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
		return nil, fmt.Errorf("Report is empty")
	}

	reports = rc.Cleaner.RemoveEmptyRows(reports)
	if len(reports) == 0 {
		return nil, fmt.Errorf("Removing empty rows resulted in empty report")
	}
	reports = rc.Cleaner.RemoveShortAndTruncateLongRows(reports)
	if len(reports) == 0 {
		return nil, fmt.Errorf("Removing short rows resulted in empty report")
	}

	return reports, nil
}
