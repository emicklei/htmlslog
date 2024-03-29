package htmlslog

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
)

type Options struct {
	Title string
	Level slog.Level
}

type Handler struct {
	buf     *bytes.Buffer
	attrs   []slog.Attr
	options Options
}

// WithAttrs implements slog.Handler.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h.attrs = append(h.attrs, attrs...)
	return h
}

// WithLevel implements slog.Handler.
func (h *Handler) WithLevel(level slog.Level) slog.Handler {
	h.options.Level = level
	return h
}

// WithGroup implements slog.Handler.
func (h *Handler) WithGroup(name string) slog.Handler {
	return h
}

func New(options Options) *Handler {
	h := &Handler{
		buf:     new(bytes.Buffer),
		options: options,
	}
	h.beginHTML()
	return h
}
func (h *Handler) beginHTML() {
	h.buf.WriteString(`
	<!DOCTYPE html>
	<head><title>Log</title></head>
	<style>
	body {
		font-family: "Open Sans", Tahoma, Geneva, sans-serif;
		-webkit-font-smoothing: antialiased;
		-moz-osx-font-smoothing: grayscale;
	}	
	</style>
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
	return h.options.Level <= l
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
	for _, a := range h.attrs {
		h.writeAttr(a.Key, a.Value)
	}
	rec.Attrs(func(a slog.Attr) bool {
		h.writeAttr(a.Key, a.Value)
		return true
	})
	h.buf.WriteString("</td>")

	h.buf.WriteString("</tr>")
	return nil
}
func (h *Handler) writeAttr(key string, value any) {
	h.buf.WriteString(key)
	h.buf.WriteString("=<b>")
	fmt.Fprintf(h.buf, "%v", value)
	h.buf.WriteString("</b> ")
}
