package cli

import (
	"fmt"
	"os"
	"time"
)

// saveToFile writes the provided data to a file with the specified filename.
// It creates the file if it doesn't exist or overwrites it if it does.
// The file is created with permissions 0644 (readable by everyone, writable by owner).
//
// Parameters:
//   - filename: The path to the file where data should be written
//   - data: The string content to write to the file
//
// Returns:
//   - error: nil if successful, otherwise an error describing what went wrong
func saveToFile(filename, data string) error {
	return os.WriteFile(filename, []byte(data), 0644)
}

// startSpinner displays an animated spinner in the console to indicate ongoing activity.
// It runs in a loop until a signal is received on the done channel.
// This provides visual feedback to users during long-running operations.
//
// Parameters:
//   - message: The text message to display alongside the spinner
//   - done: A channel that, when closed or receives a value, stops the spinner and displays completion
//
// The function blocks until the done channel receives a value or is closed.
func startSpinner(message string, done chan bool) {
	spinChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0
	for {
		select {
		case <-done:
			fmt.Printf("\r%s... Done!     \n", message)
			return
		default:
			fmt.Printf("\r%s %s", spinChars[i], message)
			time.Sleep(100 * time.Millisecond)
			i = (i + 1) % len(spinChars)
		}
	}
}
