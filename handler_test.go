package htmlslog

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"testing"
)

func TestHandle(t *testing.T) {
	h := New(Options{Level: slog.LevelInfo, Title: "Test"})
	base := slog.New(h)
	l := base.With("ctx", "summer")
	l.Debug("test", "attr", "values")
	l.Info("info", "attr", "values")
	l.Error("error", "err", errors.New("test error"), "why", "because")
	l.Info("info", "html", "<a href=\"http://example.com\">html not allowed</a>")

	o, _ := os.Create("handler_test.html")
	io.WriteString(o, h.Close())
	o.Close()
}
