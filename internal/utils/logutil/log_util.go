package logutil

import (
	"os"
	"strings"

	"github.com/phuslu/log"
)

var allLevels = []string{
	log.TraceLevel.String(),
	log.DebugLevel.String(),
	log.InfoLevel.String(),
	log.WarnLevel.String(),
	log.ErrorLevel.String(),
	log.FatalLevel.String(),
	log.PanicLevel.String(),
}

func ConfigureDefaultLogger() {
	if !log.IsTerminal(os.Stderr.Fd()) {
		return
	}
	log.DefaultLogger = log.Logger{
		TimeFormat: "15:04:05",
		Caller:     1,
		Writer: &log.ConsoleWriter{
			ColorOutput:    true,
			QuoteString:    true,
			EndWithMessage: true,
		},
	}
}

func SetDefaultLogLevel(logLevel string) {
	for _, level := range allLevels {
		if logLevel == level {
			log.DefaultLogger.Level = log.ParseLevel(level)
		}
	}
}

func AllLevelsAsString() string {
	return strings.Join(allLevels, ",")
}
