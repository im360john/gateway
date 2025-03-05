package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
)

const logDirName = "centralmind"

func getDefaultLogDir() string {
	if envDir := os.Getenv("GATEWAY_LOG_DIR"); envDir != "" {
		if ensureDirectoryExists(envDir) {
			logrus.Debugf("Using log directory from environment: %s", envDir)
			return envDir
		}
		logrus.Warnf("Environment-specified log directory %s is not accessible, falling back to defaults", envDir)
	}

	var defaultDir string

	switch runtime.GOOS {
	case "windows":
		// On Windows, try ProgramData first; if not available, fall back to the executable directory.
		if pd := os.Getenv("ProgramData"); pd != "" {
			defaultDir = filepath.Join(pd, logDirName)
		} else if exe, err := os.Executable(); err == nil {
			defaultDir = filepath.Join(filepath.Dir(exe), logDirName)
		} else {
			logrus.WithError(err).Warn("Unable to get executable path, using current directory as fallback")
			defaultDir = filepath.Join(".", logDirName)
		}
	case "darwin": // macOS
		// On macOS, default to a hidden directory inside the executable's folder.
		if exe, err := os.Executable(); err == nil {
			defaultDir = filepath.Join(filepath.Dir(exe), fmt.Sprintf(".%s", logDirName))
		} else {
			logrus.WithError(err).Warn("Unable to get executable path, using current directory as fallback")
			defaultDir = filepath.Join(".", logDirName)
		}
	default: // Linux and other Unix-like systems
		defaultDir = filepath.Join("/var/log", logDirName)
		// If creating /var/log/centralmind fails, fall back to a hidden directory in the home folder.
		if !ensureDirectoryExists(defaultDir) {
			if homeDir, err := os.UserHomeDir(); err == nil {
				defaultDir = filepath.Join(homeDir, fmt.Sprintf(".%s", logDirName))
			} else {
				logrus.WithError(err).Warn("Unable to determine user home directory, using current directory as fallback")
				defaultDir = filepath.Join(".", logDirName)
			}
		}
	}

	// Ensure the chosen directory exists.
	if !ensureDirectoryExists(defaultDir) {
		logrus.Warnf("Failed to create log directory %s, using current directory as last resort", defaultDir)
		defaultDir = "."
	}

	logrus.Debugf("Using log directory: %s", defaultDir)
	return defaultDir
}

// ensureDirectoryExists checks if a directory exists and creates it if necessary.
// It returns true if the directory exists or was created successfully.
func ensureDirectoryExists(dir string) bool {
	// Check if the directory already exists.
	if fi, err := os.Stat(dir); err == nil && fi.IsDir() {
		return true
	}

	// Try to create the directory.
	if err := os.MkdirAll(dir, 0755); err != nil {
		logrus.WithError(err).WithField("directory", dir).Warn("Failed to create log directory")
		return false
	}

	logrus.WithField("directory", dir).Info("Created log directory")
	return true
}
