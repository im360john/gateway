package cli

import (
	"fmt"
	"os"
	"time"
)

func saveToFile(filename, data string) error {
	return os.WriteFile(filename, []byte(data), 0644)
}

// startSpinner starts a loading animation in the console
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
