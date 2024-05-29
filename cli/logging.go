package cli

import (
	"log/slog"
	"os"
)

var Log = newLog()

func newLog() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, nil))
}
