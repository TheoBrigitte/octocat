package json

import (
	"encoding/json"
	"io"

	log "github.com/sirupsen/logrus"
)

// Config used to initialize a new JSON storage.
type Config struct {
	Writer io.Writer
	Logger *log.Logger
}

// JSON is a superset of json.Encoder.
// It stores data in a json file.
type JSON struct {
	logger *log.Logger

	*json.Encoder
}

// New creates a new JSON.
func New(config Config) JSON {
	return JSON{
		Encoder: json.NewEncoder(config.Writer),
		logger:  config.Logger,
	}
}
