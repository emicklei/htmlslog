package htmlslog

import (
	"errors"
	"log/slog"
	"os"
	"testing"
)

func TestHandle(t *testing.T) {
	o, _ := os.Create("handler_test.html")
	defer o.Close()

	h := New(o, Options{Level: slog.LevelInfo, Title: "Test"})
	defer h.Close()

	base := slog.New(h)
	l := base.With("ctx", "summer")
	l.Debug("test", "attr", "values")
	l.Info("info", "attr", "values")
	l.Error("error", "err", errors.New("test error"), "why", "because")
	l.Info("info", "html", "<a href=\"http://example.com\">html not allowed</a>")
}
