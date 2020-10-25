// Package util provides some utilities
package util

import (
	"os"

	"github.com/comail/colog"
)

// InitializeLog initialize log settings
func InitializeLog(verbose bool) {
	colog.Register()
	colog.SetOutput(os.Stderr)
	colog.SetDefaultLevel(colog.LDebug)
	colog.SetMinLevel(colog.LInfo)

	if verbose {
		colog.SetMinLevel(colog.LDebug)
	}
}
