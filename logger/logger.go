package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Event stores messages to log later, from our standard interface
type Event struct {
	id      int
	message string
}

// StandardLogger enforces specific log message formats
type StandardLogger struct {
	baseLogger zerolog.Logger
}

// NewLogger initializes the standard logger
func NewLogger(verbosity int) *StandardLogger {
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC822}
	baseLogger := newBaseLogger(verbosity, consoleWriter)
	var standardLogger = &StandardLogger{baseLogger}

	return standardLogger
}

func NewLoggerWithWriter(verbosity int, out io.Writer) *StandardLogger {
	consoleWriter := zerolog.ConsoleWriter{Out: out, TimeFormat: time.RFC822}
	baseLogger := newBaseLogger(verbosity, consoleWriter)
	var standardLogger = &StandardLogger{baseLogger}

	return standardLogger
}
func newBaseLogger(verbosity int, consoleWriter zerolog.ConsoleWriter) zerolog.Logger {
	var baseLogger = zerolog.New(consoleWriter).With().Timestamp().Logger()
	baseLogger = baseLogger.Level(zerolog.ErrorLevel)
	switch verbosity {
	case 0:
		baseLogger = baseLogger.Level(zerolog.ErrorLevel)
	case 1:
		baseLogger = baseLogger.Level(zerolog.WarnLevel)
	case 2:
		baseLogger = baseLogger.Level(zerolog.InfoLevel)
	case 3:
		baseLogger = baseLogger.Level(zerolog.DebugLevel)
	case 4:
		baseLogger = baseLogger.Level(zerolog.TraceLevel)
	}

	return baseLogger
}

// Declare variables to store log messages as new Events
var (
	invalidArgMessage      = Event{1, "Invalid arg: %s"}
	invalidArgValueMessage = Event{2, "Invalid value for argument: %s: %v"}
	missingArgMessage      = Event{3, "Missing arg: %s"}
)

// InvalidArg is a standard error message
func (l *StandardLogger) InvalidArgs(messag string, arguments []string) {
	l.baseLogger.Error().Str("status", "debug").Msgf(invalidArgMessage.message, arguments)

}

// InvalidArgValue is a standard error message
func (l *StandardLogger) InvalidArgValue(argumentName string, argumentValue string) {
	l.baseLogger.Error().Str(
		"status",
		"error",
	).Msgf(
		invalidArgValueMessage.message,
		argumentName,
		argumentValue,
	)
}

// MissingArg is a standard error message
func (l *StandardLogger) MissingArg(argumentName string) {
	l.baseLogger.Error().Str("status", "error").Msgf(missingArgMessage.message, argumentName)
}

func (l *StandardLogger) Debugf(message string, args ...interface{}) {
	l.baseLogger.Debug().Str("status", "debug").Msgf(message, args...)
}

func (l *StandardLogger) Debug(message string) {
	l.baseLogger.Debug().Str("status", "debug").Msg(message)
}

func (l *StandardLogger) Infof(message string, args ...interface{}) {
	l.baseLogger.Info().Str("status", "info").Msgf(message, args...)
}

func (l *StandardLogger) Info(message string) {
	l.baseLogger.Info().Str("status", "info").Msg(message)
}

func (l *StandardLogger) Error(err error, message string) {
	l.baseLogger.Error().Str("status", "error").Err(err).Msg(message)
}

func (l *StandardLogger) Errorf(err error, message string, args ...interface{}) {
	l.baseLogger.Error().Str("status", "error").Err(err).Msgf(message, args)
}

func (l *StandardLogger) Fatal(err error, message string) {
	l.baseLogger.Fatal().Str("status", "fatal").Err(err).Msg(message)
}

func (l *StandardLogger) Fatalf(err error, message string, args ...interface{}) {
	l.baseLogger.Fatal().Str("status", "fatal").Err(err).Msgf(message, args)
}

func (l *StandardLogger) WithContext(key, value string) {
	l.baseLogger = l.baseLogger.With().Str(key, value).Logger()
}
