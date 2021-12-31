// Package terminal deals with terminal operations like progress and profiles
package terminal

import (
	"log"
	"time"

	"github.com/briandowns/spinner"
)

var (
	spinnerDelay   = 100 * time.Millisecond
	spinnerCharset = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	spinnerColor   = "blue"
)

type Spinner struct {
	*spinner.Spinner
}

func NewSpinner(message string) *Spinner {
	// create basic spinner coloured spinner
	s := spinner.New(spinnerCharset, spinnerDelay, spinner.WithWriter(log.Writer()))
	s.Color(spinnerColor)
	s.Suffix = " " + message

	// clear terminal line after spinner is stopped
	s.FinalMSG = "\033[2K\r"

	return &Spinner{s}
}

func WithSpinner(message string, f func() error) error {
	spinner := NewSpinner(message)
	spinner.Start()
	err := f()
	spinner.Stop()
	return err
}
