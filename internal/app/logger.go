package app

import (
	"time"

	"github.com/rs/zerolog"
)

func NewLogger() *zerolog.Logger {
	out := zerolog.NewConsoleWriter(
		func(w *zerolog.ConsoleWriter) {
			w.TimeFormat = time.RFC3339
		},
	)
	log := zerolog.New(out).With().Timestamp().Logger()
	return &log
}
