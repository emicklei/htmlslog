package htmlslog

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"html"
	"log/slog"
	"strings"
)

// TODO add io.Writer
type Options struct {
	Title string
	Level slog.Level
}

type Handler struct {
	buf     *bytes.Buffer
	odd     bool // for highlighting rows
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

//go:embed prolog.html
var prolog string

func (h *Handler) beginHTML() {
	if h.options.Title != "" {
		prolog = strings.ReplaceAll(prolog, "TITLE", h.options.Title)
	}
	h.buf.WriteString(prolog)
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
	if h.odd {
		h.buf.WriteString("<tr class=\"odd\">")
	} else {
		h.buf.WriteString("<tr>")
	}
	h.odd = !h.odd

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
	var valueString string
	switch value.(type) {
	case string:
		valueString = value.(string)
	case error:
		valueString = value.(error).Error()
	default:
		valueString = fmt.Sprintf("%v", value)
	}
	fmt.Fprintf(h.buf, "%v", html.EscapeString(valueString))
	h.buf.WriteString("</b> ")
}
