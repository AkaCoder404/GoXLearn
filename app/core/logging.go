package core

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/cloudwego/kitex/pkg/utils"
)

const (
	dateFormat    = "2006-01-02 15-04-05"
	hoursPerMonth = float64(24 * 31)
)

// loggingDir: The logging directory to store the logs

var loggingDir = filepath.Join(utils.GetConfDir(), "logs")

// setUpLogging: Set up the logger to log ay useful information when application is running
func (m *GoLearn) setUpLogging() error {
	if err := os.MkdirAll(loggingDir, os.ModePerm); err != nil {
		return err
	}

	// Remove old logging files (at least one-month-old)
	now := time.Now()
	_ = filepath.Walk(loggingDir, func(path string, info fs.FileInfo, err error) error {
		// Remove files that were modified more than 1 month ago.
		fileDate := info.ModTime()
		// Ignore folders, even though there should not be any in this folder.
		if !info.IsDir() && now.Sub(fileDate).Hours() >= hoursPerMonth {
			_ = os.Remove(path)
		}
		return nil
	})

	// Create file for current session logging
	formattedDate := now.Format(dateFormat)
	logFilePath := filepath.Join(loggingDir, fmt.Sprintf("%s.log", formattedDate))

	var err error
	if m.LogFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm); err != nil {
		return err
	}
	log.SetOutput(m.LogFile)
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Printf("Session started at %s\n", formattedDate)

	return nil
}

// stopLogging : Closes the log file
func (m *GoLearn) stopLogging() error {
	return m.LogFile.Close()
}
