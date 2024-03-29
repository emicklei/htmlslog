package htmlslog

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"testing"
)

func TestHandle(t *testing.T) {
	h := New(slog.LevelInfo)
	l := slog.New(h)
	l.Debug("test", "attr", "values")
	l.Info("info", "attr", "values")
	l.Error("error", "err", errors.New("test error"), "why", "because")

	o, _ := os.Create("handler_test.html")
	io.WriteString(o, h.Close())
	o.Close()
}
