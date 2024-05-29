package cli

import (
	"log/slog"
	"os"
)

func Log() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, nil))
}
