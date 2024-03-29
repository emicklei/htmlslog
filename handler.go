package htmlslog

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
)

type Handler struct {
	buf   *bytes.Buffer
	level slog.Level
	attrs []slog.Attr
}

// WithAttrs implements slog.Handler.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h.attrs = append(h.attrs, attrs...)
	return h
}

// WithLevel implements slog.Handler.
func (h *Handler) WithLevel(level slog.Level) slog.Handler {
	h.level = level
	return h
}

// WithGroup implements slog.Handler.
func (h *Handler) WithGroup(name string) slog.Handler {
	return h
}

func New(minLevel slog.Level) *Handler {
	h := &Handler{
		buf:   new(bytes.Buffer),
		level: minLevel,
	}
	h.beginHTML()
	return h
}
func (h *Handler) beginHTML() {
	h.buf.WriteString(`
	<!DOCTYPE html>
	<head><title>Log</title></head>
	<html><body>
		<table>
		<tr>
			<th>Time</th>	
			<th>Level</th>
			<th>Message</th>
			<th>Attrs</th>
		</tr>
`)
}
func (h *Handler) Close() string {
	h.endHTML()
	return h.buf.String()
}
func (h *Handler) endHTML() {
	h.buf.WriteString("</table></body></html>")
}
func (h *Handler) Enabled(ctx context.Context, l slog.Level) bool {
	return h.level <= l
}
func (h *Handler) Handle(ctx context.Context, rec slog.Record) error {
	h.buf.WriteString("<tr>")

	h.buf.WriteString("<td>")
	h.buf.WriteString(rec.Time.Format("2006-01-02 15:04:05"))
	h.buf.WriteString("</td>")

	h.buf.WriteString("<td>")
	h.buf.WriteString(rec.Level.String())
	h.buf.WriteString("</td>")

	h.buf.WriteString("<td>")
	h.buf.WriteString(rec.Message)
	h.buf.WriteString("</td>")

	h.buf.WriteString("<td>")
	rec.Attrs(func(a slog.Attr) bool {
		h.buf.WriteString(a.Key)
		h.buf.WriteString("=")
		fmt.Fprintf(h.buf, "%v", a.Value)
		h.buf.WriteString(" ")
		return true
	})
	h.buf.WriteString("</td>")

	h.buf.WriteString("</tr>")
	return nil
}
