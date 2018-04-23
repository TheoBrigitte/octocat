package csv

import (
	"encoding/csv"
	"io"

	log "github.com/sirupsen/logrus"
)

// Config used to initialize a new CSV storage.
type Config struct {
	Writer io.Writer
	Header []string
	Logger *log.Logger
}

// CSV is a superset of csv.Writer.
// It stores data in a csv file.
type CSV struct {
	logger *log.Logger
	header []string

	*csv.Writer
}

// Formatter is an interface used to retrive CSV columns.
type Formatter interface {
	CSV() []string
}

// New creates a new CSV.
func New(config Config) CSV {
	return CSV{
		Writer: csv.NewWriter(config.Writer),
		header: config.Header,
		logger: config.Logger,
	}
}

// Flush data to the file.
func (s CSV) Flush() error {
	s.Writer.Flush()

	return s.Writer.Error()
}
