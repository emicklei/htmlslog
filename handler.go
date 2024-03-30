package htmlslog

import (
	"context"
	_ "embed"
	"fmt"
	"html"
	"io"
	"log/slog"
	"strings"
)

type Options struct {
	Title              string
	Level              slog.Level
	PassthroughHandler slog.Handler // if set then also handle events by this handler
}

type Handler struct {
	writer  io.Writer
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

func New(w io.Writer, options Options) *Handler {
	h := &Handler{options: options, writer: w}
	h.beginHTML()
	return h
}

//go:embed prolog.html
var prolog string

func (h *Handler) beginHTML() {
	if h.options.Title != "" {
		prolog = strings.ReplaceAll(prolog, "TITLE", h.options.Title)
	}
	fmt.Fprint(h.writer, prolog)
	fmt.Fprint(h.writer, "\n")
}

// Close write the ending HTML and return the result.
func (h *Handler) Close() {
	h.endHTML()
}
func (h *Handler) endHTML() {
	fmt.Fprint(h.writer, "</table></body></html>")
}

// Enabled implements slog.Handler.
func (h *Handler) Enabled(ctx context.Context, l slog.Level) bool {
	return h.options.Level <= l
}

// Handle implements slog.Handler.
func (h *Handler) Handle(ctx context.Context, rec slog.Record) error {
	if h.odd {
		fmt.Fprint(h.writer, "<tr class=\"odd\">")
	} else {
		fmt.Fprint(h.writer, "<tr>")
	}
	h.odd = !h.odd

	fmt.Fprint(h.writer, "<td>")
	fmt.Fprint(h.writer, rec.Time.Format("2006-01-02 15:04:05"))
	fmt.Fprint(h.writer, "</td>")

	fmt.Fprint(h.writer, "<td>")
	fmt.Fprint(h.writer, rec.Level.String())
	fmt.Fprint(h.writer, "</td>")

	fmt.Fprint(h.writer, "<td>")
	fmt.Fprint(h.writer, rec.Message)
	fmt.Fprint(h.writer, "</td>")

	fmt.Fprint(h.writer, "<td>")
	for _, a := range h.attrs {
		h.writeAttr(a.Key, a.Value)
	}
	rec.Attrs(func(a slog.Attr) bool {
		h.writeAttr(a.Key, a.Value)
		return true
	})
	fmt.Fprint(h.writer, "</td>")

	fmt.Fprint(h.writer, "</tr>\n")

	if h.options.PassthroughHandler != nil {
		return h.options.PassthroughHandler.Handle(ctx, rec)
	}
	return nil
}
func (h *Handler) writeAttr(key string, value any) {
	fmt.Fprint(h.writer, key)
	fmt.Fprint(h.writer, "=<b>")
	var valueString string
	switch value.(type) {
	case string:
		valueString = value.(string)
	case error:
		valueString = value.(error).Error()
	default:
		valueString = fmt.Sprintf("%v", value)
	}
	fmt.Fprintf(h.writer, "%v", html.EscapeString(valueString))
	fmt.Fprint(h.writer, "</b> ")
}
