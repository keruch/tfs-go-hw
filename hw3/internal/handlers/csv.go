package handlers

import (
	"encoding/csv"
	"os"
)

// CSVRecord interface intended for working with csv
// Record creates mapping from some data (for example struct field) to CSV string
type CSVRecord interface {
	Record() []string
}

type CSVWriterFlusher interface {
	Write(record CSVRecord) error
	Flush()
}

type CSVWriter struct {
	writer *csv.Writer
}

func NewCSVWriterCloser(filename string) (CSVWriterFlusher, error) {
	file, err := os.Create(filename)
	if err != nil {
		return CSVWriter{}, err
	}

	cw := CSVWriter{
		writer: csv.NewWriter(file),
	}

	return cw, nil
}

func (cw CSVWriter) Write(record CSVRecord) error {
	err := cw.writer.Write(record.Record())
	if err != nil {
		return err
	}

	return nil
}

func (cw CSVWriter) Flush() {
	cw.writer.Flush()
}
