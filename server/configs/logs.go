package configs

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func InitLogs() log.Logger {
	return log.Logger{
		Out:          os.Stdout,
		Level:        log.DebugLevel,
		//Formatter:    &log.JSONFormatter{}, // uncomment this to enable JSON logs
	}
}
