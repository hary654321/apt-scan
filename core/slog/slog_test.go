package slog

import (
	"testing"
)

func TestName(t *testing.T) {
	Printf(INFO, "INFO")
	Printf(WARN, "WARN")
	Printf(DEBUG, "DEBUG")
	Printf(NONE, "NONE")
	Printf(ERROR, "ERROR")
}
