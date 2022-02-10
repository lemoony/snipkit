package logutil

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/phuslu/log"

	"github.com/lemoony/snipkit/internal/utils/system"
)

const maxBackups = 2

var allLevels = []string{
	log.TraceLevel.String(),
	log.DebugLevel.String(),
	log.InfoLevel.String(),
	log.WarnLevel.String(),
	log.ErrorLevel.String(),
	log.FatalLevel.String(),
	log.PanicLevel.String(),
}

func ConfigureDefaultLogger(s *system.System) {
	if !log.IsTerminal(os.Stderr.Fd()) {
		return
	}

	fileWriter := log.FileWriter{
		Filename:     filepath.Join(s.HomeDir(), ".log", "log"),
		EnsureFolder: true,
		MaxBackups:   maxBackups,
	}

	log.DefaultLogger = log.Logger{
		TimeFormat: "15:04:05",
		Caller:     1,
		Writer: &log.ConsoleWriter{
			QuoteString:    true,
			EndWithMessage: true,
			Writer:         &fileWriter,
		},
	}

	_ = fileWriter.Rotate()
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
