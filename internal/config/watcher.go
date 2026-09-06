package config

import (
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Watcher monitors a file for changes using polling (cross-platform).
type Watcher struct {
	filePath  string
	onChange  func()
	stopCh    chan struct{}
	stoppedCh chan struct{}
}

// NewWatcher creates a new file watcher that polls the file for modifications.
func NewWatcher(filePath string, onChange func()) (*Watcher, error) {
	w := &Watcher{
		filePath:  filePath,
		onChange:  onChange,
		stopCh:    make(chan struct{}),
		stoppedCh: make(chan struct{}),
	}

	go w.watch()
	return w, nil
}

func (w *Watcher) watch() {
	defer close(w.stoppedCh)

	var lastModTime time.Time
	info, err := os.Stat(w.filePath)
	if err == nil {
		lastModTime = info.ModTime()
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			info, err := os.Stat(w.filePath)
			if err != nil {
				continue
			}
			if info.ModTime().After(lastModTime) {
				lastModTime = info.ModTime()
				slog.Debug("config file changed", "path", w.filePath)
				w.onChange()
			}
		case <-w.stopCh:
			return
		}
	}
}

// Stop stops the file watcher.
func (w *Watcher) Stop() {
	close(w.stopCh)
	<-w.stoppedCh
}

// EnsureConfigFile creates a default config file if one doesn't exist.
func EnsureConfigFile(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	cfg := DefaultConfig()
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}