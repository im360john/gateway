package cli

import (
	"os"
	"path/filepath"
)

// getDefaultLogDir returns the default log directory based on the operating system
func getDefaultLogDir() string {
	if os.PathSeparator == '\\' { // Windows
		// Get the executable's directory
		exe, err := os.Executable()
		if err != nil {
			return "."
		}
		return filepath.Dir(exe)
	}
	// Unix-like systems
	return "/var/log/gateway"
}
