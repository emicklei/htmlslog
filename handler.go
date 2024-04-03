package htmlslog

import (
	"context"
	_ "embed"
	"fmt"
	"html"
	"io"
	"log/slog"
	"strings"
	"time"
)

type Options struct {
	Title              string
	TimeLayout         string
	Level              slog.Level
	TableOnly          bool         // if true then only the table and its rows are written
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
	if h.options.PassthroughHandler != nil {
		h.options.PassthroughHandler = h.options.PassthroughHandler.WithAttrs(attrs)
	}
	return h
}

// WithLevel implements slog.Handler.
func (h *Handler) WithLevel(level slog.Level) slog.Handler {
	h.options.Level = level
	return h
}

// WithGroup implements slog.Handler.
func (h *Handler) WithGroup(name string) slog.Handler {
	if h.options.PassthroughHandler != nil {
		h.options.PassthroughHandler = h.options.PassthroughHandler.WithGroup(name)
	}
	return h
}

func New(w io.Writer, options Options) *Handler {
	if options.TimeLayout == "" {
		options.TimeLayout = time.RFC3339
	}
	if options.Level == 0 {
		options.Level = slog.LevelInfo
	}
	h := &Handler{options: options, writer: w}
	h.beginHTML()
	return h
}

//go:embed prolog.html
var prolog string

func (h *Handler) beginHTML() {
	if !h.options.TableOnly {
		if h.options.Title != "" {
			prolog = strings.ReplaceAll(prolog, "TITLE", h.options.Title)
		}
		fmt.Fprint(h.writer, prolog)
		fmt.Fprint(h.writer, "\n")
	}
	h.beginTable()
}

func (h *Handler) beginTable() {
	fmt.Fprint(h.writer, `<table>
<tr>
	<th>Time</th>
	<th>Level</th>
	<th>Message</th>
	<th>Attributes</th>
</tr>
`)
}

// Close write the ending HTML and return the result.
func (h *Handler) Close() {
	h.endHTML()
}
func (h *Handler) endTable() {
	fmt.Fprint(h.writer, "</table>")
}
func (h *Handler) endHTML() {
	h.endTable()
	if !h.options.TableOnly {
		fmt.Fprint(h.writer, "</body></html>")
	}
}

// Enabled implements slog.Handler.
func (h *Handler) Enabled(ctx context.Context, l slog.Level) bool {
	return h.options.Level <= l
}

// Handle implements slog.Handler.
func (h *Handler) Handle(ctx context.Context, rec slog.Record) error {
	if h.odd {
		fmt.Fprint(h.writer, "<tr class=\"odd\">\n")
	} else {
		fmt.Fprint(h.writer, "<tr>\n")
	}
	h.odd = !h.odd

	fmt.Fprint(h.writer, "\t<td>")
	fmt.Fprint(h.writer, rec.Time.Format(h.options.TimeLayout))
	fmt.Fprint(h.writer, "</td>\n")

	fmt.Fprint(h.writer, "\t<td>")
	fmt.Fprint(h.writer, rec.Level.String())
	fmt.Fprint(h.writer, "</td>\n")

	fmt.Fprint(h.writer, "\t<td>")
	fmt.Fprint(h.writer, rec.Message)
	fmt.Fprint(h.writer, "</td>\n")

	fmt.Fprint(h.writer, "\t<td>")
	for _, a := range h.attrs {
		h.writeAttr(a.Key, a.Value)
	}
	rec.Attrs(func(a slog.Attr) bool {
		h.writeAttr(a.Key, a.Value)
		return true
	})
	fmt.Fprint(h.writer, "</td>\n")

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
