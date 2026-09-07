// Package webui provides the embedded web management interface.
package webui

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

// logEntry represents a single log entry.
type logEntry struct {
	Time    time.Time `json:"time"`
	Level   string    `json:"level"`
	Message string    `json:"message"`
}

// logBuffer is a circular buffer that stores recent log entries
// and allows SSE streaming to connected clients.
type logBuffer struct {
	mu       sync.RWMutex
	entries  []logEntry
	maxSize  int
	subs     map[int]chan string
	subID    int
}

// NewLogBuffer creates a new log buffer with the given max size.
func NewLogBuffer(maxSize int) *logBuffer {
	return &logBuffer{
		entries: make([]logEntry, 0, maxSize),
		maxSize: maxSize,
		subs:    make(map[int]chan string),
	}
}

// append adds a new log entry to the buffer.
func (lb *logBuffer) append(level string, msg string) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	entry := logEntry{
		Time:    time.Now(),
		Level:   level,
		Message: msg,
	}

	lb.entries = append(lb.entries, entry)
	if len(lb.entries) > lb.maxSize {
		lb.entries = lb.entries[len(lb.entries)-lb.maxSize:]
	}

	// Format as a simple text line for SSE
	line := formatLogLine(entry)

	// Send to all subscribers (non-blocking)
	for _, ch := range lb.subs {
		select {
		case ch <- line:
		default:
			// Drop if subscriber is too slow
		}
	}
}

// subscribe creates a new subscription channel for SSE streaming.
// The channel receives formatted log lines.
func (lb *logBuffer) subscribe(ctx context.Context) <-chan string {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	ch := make(chan string, 64)
	id := lb.subID
	lb.subID++
	lb.subs[id] = ch

	// Send existing entries
	for _, entry := range lb.entries {
		ch <- formatLogLine(entry)
	}

	// Clean up on context cancel
	go func() {
		<-ctx.Done()
		lb.mu.Lock()
		delete(lb.subs, id)
		lb.mu.Unlock()
	}()

	return ch
}

// formatLogLine formats a log entry as a text line.
func formatLogLine(entry logEntry) string {
	return entry.Time.Format("15:04:05") + " " + entry.Level + " " + entry.Message
}

// slogHandler is a slog.Handler that forwards log records to the log buffer.
type slogHandler struct {
	next  slog.Handler
	buf   *logBuffer
	level slog.Leveler
}

// NewSlogHandler creates a new slog handler that forwards to the log buffer.
func NewSlogHandler(next slog.Handler, buf *logBuffer, level slog.Leveler) *slogHandler {
	return &slogHandler{
		next:  next,
		buf:   buf,
		level: level,
	}
}

// Enabled reports whether the handler handles records at the given level.
func (h *slogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

// Handle handles the Record.
func (h *slogHandler) Handle(ctx context.Context, r slog.Record) error {
	msg := r.Message
	r.Attrs(func(a slog.Attr) bool {
		msg += " " + a.Key + "=" + a.Value.String()
		return true
	})
	h.buf.append(r.Level.String(), msg)
	return h.next.Handle(ctx, r)
}

// WithAttrs returns a new handler with the given attributes.
func (h *slogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &slogHandler{
		next:  h.next.WithAttrs(attrs),
		buf:   h.buf,
		level: h.level,
	}
}

// WithGroup returns a new handler with the given group.
func (h *slogHandler) WithGroup(name string) slog.Handler {
	return &slogHandler{
		next:  h.next.WithGroup(name),
		buf:   h.buf,
		level: h.level,
	}
}