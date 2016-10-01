package csv

import (
	"encoding/csv"
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
		return nil, err
	}
	reports = rc.Cleaner.RemoveEmptyRows(reports)
	reports = rc.Cleaner.RemoveIrregularLengthRows(reports)
	return reports, nil
}
