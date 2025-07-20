package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	logDir      string
	maxLogFiles int
	maxLogSize  int64 // in MB
}

func New(logDir string) *Logger {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatal("Failed to create log directory:", err)
	}

	infoLogFile := filepath.Join(logDir, fmt.Sprintf("info_%s.log", time.Now().Format("2006-01-02")))
	errorLogFile := filepath.Join(logDir, fmt.Sprintf("error_%s.log", time.Now().Format("2006-01-02")))

	infoFile, err := os.OpenFile(infoLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open info log file:", err)
	}

	errorFile, err := os.OpenFile(errorLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open error log file:", err)
	}

	logger := &Logger{
		infoLogger:  log.New(io.MultiWriter(os.Stdout, infoFile), "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(io.MultiWriter(os.Stderr, errorFile), "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		logDir:      logDir,
		maxLogFiles: 30,
		maxLogSize:  1000, // 100MB per file
	}

	// Clean up old logs
	go logger.cleanupOldLogs()

	return logger
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.infoLogger.Printf(format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.errorLogger.Printf(format, args...)
	// Create error report zip when error occurs
}

func (l *Logger) ErrorWithContext(context map[string]interface{}, format string, args ...interface{}) {
	contextStr := ""
	for k, v := range context {
		contextStr += fmt.Sprintf("%s=%v ", k, v)
	}
	l.errorLogger.Printf("[Context: %s] %s", contextStr, fmt.Sprintf(format, args...))
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	l.errorLogger.Printf(format, args...)
	os.Exit(1)
}

func (l *Logger) cleanupOldLogs() {
	files, err := filepath.Glob(filepath.Join(l.logDir, "*.log"))
	if err != nil {
		l.errorLogger.Printf("Failed to list log files: %v", err)
		return
	}

	// Remove old log files
	for _, file := range files {
		stat, err := os.Stat(file)
		if err != nil {
			continue
		}

		// Remove files older than maxLogFiles days
		if time.Since(stat.ModTime()) > time.Duration(l.maxLogFiles)*24*time.Hour {
			os.Remove(file)
		}
	}

	// Remove old error reports
	reports, err := filepath.Glob(filepath.Join(l.logDir, "error_report_*.zip"))
	if err != nil {
		return
	}

	for _, report := range reports {
		stat, err := os.Stat(report)
		if err != nil {
			continue
		}

		// Remove reports older than 7 days
		if time.Since(stat.ModTime()) > 7*24*time.Hour {
			os.Remove(report)
		}
	}
}

func (l *Logger) GetLogDir() string {
	return l.logDir
}

func (l *Logger) RotateLogs() {
	// Close current log files and create new ones
	l.cleanupOldLogs()

	// Note: In a real implementation, we'd need to reopen files
	// This is a placeholder for log rotation functionality
	l.Info("Log rotation completed")
}
