package i2c

import (
	"os"

	"github.com/op/go-logging"
)

// Comment INFO and uncomment DEBUG if you want detail debug output in library.
var log *logging.Logger = buildLogger("i2c",
	//	logging.DEBUG,
	logging.INFO,
)

var terminalBackend logging.LeveledBackend = nil

func buildLogger(module string, level logging.Level) *logging.Logger {
	// Set the backends to be used.
	if terminalBackend == nil {
		// Everything except the message has a custom color
		// which is dependent on the log level. Many fields have a custom output
		// formatting too, eg. the time returns the hour down to the milli second.
		var format = logging.MustStringFormatter(
			"%{time:2006-01-02T15:04:05.000} [%{module}] " +
				"%{color}%{level:.4s}%{color:reset}  %{message}",
		)
		// Create backend for os.Stderr.
		var backend logging.Backend = logging.NewLogBackend(os.Stderr, "", 0)

		// For messages written to backend we want to add some additional
		// information to the output, including the used log level and the name of
		// the function.
		var backendFormatter logging.Backend = logging.NewBackendFormatter(backend, format)
		var backendLeveled logging.LeveledBackend = logging.AddModuleLevel(backendFormatter)
		terminalBackend = backendLeveled
		logging.SetBackend(terminalBackend)
	}
	log := logging.MustGetLogger(module)
	terminalBackend.SetLevel(level, module)
	return log
}
