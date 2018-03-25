package i2c

import logger "github.com/d2r2/go-logger"

// You can manage verbosity of log output
// in the package by changing last parameter value.
var lg = logger.NewPackageLogger("i2c",
	logger.DebugLevel,
	// logger.InfoLevel,
)
