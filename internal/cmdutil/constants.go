package cmdutil

import (
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

var (
	SpinnerCharSet = spinner.CharSets[78]
	SuccessIcon    = color.GreenString("✔")
)

const (
	SpinnerInterval = 100 * time.Millisecond
	SpinnerColor    = "green"
)
