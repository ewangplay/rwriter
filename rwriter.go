package rwriter

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// LogFileMaxSize defines the default max size of each log file
const (
	LogFileBaseName = "test"
	LogFileMaxSize  = 100
	LogFilePath     = "."
)

// Config defines RotateWriter config
// Module defines the basename of log files, the defualt value is 'test';
// Path defines the path of log files, the default value is '.';
// MaxSize defines the max size of each log file, the unit is megabyte, the default value is 100;
// RotateDaily defines whether to rotate each day, the default value is false.
type Config struct {
	Module      string
	Path        string
	MaxSize     int64
	RotateDaily bool
}

// RotateWriter defines the rotate writer
type RotateWriter struct {
	cfg      *Config
	lock     sync.Mutex
	filename string
	currDate string
	fp       *os.File
	quit     chan int
}

// NewRotateWriter make a new RotateWriter. Return nil if error occurs during setup.
func NewRotateWriter(cfg *Config) (*RotateWriter, error) {

	cfg = fulfilConfig(cfg)

	w := &RotateWriter{cfg: cfg}

	w.filename = filepath.Join(cfg.Path, fmt.Sprintf("%s.log", cfg.Module))
	err := w.rotate()
	if err != nil {
		return nil, err
	}

	w.currDate = time.Now().Format("2006-01-02")
	w.quit = make(chan int)
	go w.autoRotate(w.quit)

	return w, nil
}

func fulfilConfig(cfg *Config) *Config {
	if cfg == nil {
		cfg = &Config{
			Module:      LogFileBaseName,
			Path:        LogFilePath,
			MaxSize:     LogFileMaxSize,
			RotateDaily: false,
		}
		return cfg
	}

	if cfg.Module == "" {
		cfg.Module = LogFileBaseName
	}
	if cfg.Path == "" {
		cfg.Path = LogFilePath
	}
	if cfg.MaxSize == 0 {
		cfg.MaxSize = LogFileMaxSize
	}
	return cfg
}

// Close the rotate writer
func (w *RotateWriter) Close() error {
	if w.quit != nil {
		close(w.quit)
	}
	if w.fp != nil {
		return w.fp.Close()
	}
	return nil
}

// Write satisfies the io.Writer interface.
func (w *RotateWriter) Write(output []byte) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	return w.fp.Write(output)
}

// Perform the actual act of rotating and reopening file.
func (w *RotateWriter) rotate() (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	// Close existing file if open
	if w.fp != nil {
		err = w.fp.Close()
		w.fp = nil
		if err != nil {
			return
		}
	}

	// Rename dest file if it already exists
	fileinfo, err := os.Stat(w.filename)
	if err == nil {
		if fileinfo.Size() > 0 {
			backupFilename := w.filename + "." + time.Now().Format("2006-01-02_15:04:05")
			err = os.Rename(w.filename, backupFilename)
			if err != nil {
				return
			}
		}
	}

	// Create a file.
	w.fp, err = os.Create(w.filename)
	return
}

func (w *RotateWriter) autoRotate(quit chan int) {
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-quit:
			fmt.Printf("quit auto rotate log file\n")
			return
		case <-ticker.C:

			//check log file size
			fileinfo, err := os.Stat(w.filename)
			if err == nil {
				if fileinfo.Size() >= w.cfg.MaxSize*1024*1024 {
					//rotate log file
					fmt.Printf("start to rotate log file...\n")
					err = w.rotate()
					if err != nil {
						fmt.Printf("rotate log file fail: %v\n", err)
					}
					continue
				}
			}

			//check date
			if w.cfg.RotateDaily {
				date := time.Now().Format("2006-01-02")
				if date != w.currDate {
					if fileinfo.Size() > 0 {
						//rotate log file
						fmt.Printf("start to rotate log file...\n")
						err = w.rotate()
						if err != nil {
							fmt.Printf("rotate log file fail: %v\n", err)
						} else {
							w.currDate = date
						}
						continue
					}
				}
			}
		}
	}
}
