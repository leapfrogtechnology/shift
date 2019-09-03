package spinner

import (
	"time"

	"github.com/briandowns/spinner"
)

var s = spinner.New(spinner.CharSets[14], 100*time.Millisecond)

// Start starts the spinner.
func Start(msg string) {
	s.Color("cyan")
	s.Prefix = msg

	s.Start()
}

// Stop stops the spinner.
func Stop() {
	s.Stop()
}
