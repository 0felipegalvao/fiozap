package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

func Init(level string, pretty bool) {
	var output io.Writer = os.Stdout

	if pretty {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	}

	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}

	log = zerolog.New(output).
		Level(lvl).
		With().
		Timestamp().
		Caller().
		Logger()
}

func Debug(msg string) {
	log.Debug().Msg(msg)
}

func Debugf(format string, v ...interface{}) {
	log.Debug().Msgf(format, v...)
}

func Info(msg string) {
	log.Info().Msg(msg)
}

func Infof(format string, v ...interface{}) {
	log.Info().Msgf(format, v...)
}

func Warn(msg string) {
	log.Warn().Msg(msg)
}

func Warnf(format string, v ...interface{}) {
	log.Warn().Msgf(format, v...)
}

func Error(msg string) {
	log.Error().Msg(msg)
}

func Errorf(format string, v ...interface{}) {
	log.Error().Msgf(format, v...)
}

func Fatal(msg string) {
	log.Fatal().Msg(msg)
}

func Fatalf(format string, v ...interface{}) {
	log.Fatal().Msgf(format, v...)
}

func WithField(key string, value interface{}) *zerolog.Event {
	return log.Info().Interface(key, value)
}

func WithFields(fields map[string]interface{}) *zerolog.Event {
	event := log.Info()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	return event
}

func WithError(err error) *zerolog.Event {
	return log.Error().Err(err)
}

func Get() *zerolog.Logger {
	return &log
}

func Sub(module string) zerolog.Logger {
	return log.With().Str("module", module).Logger()
}

func Writer() io.Writer {
	return os.Stdout
}
