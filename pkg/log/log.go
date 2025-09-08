package log

import (
	"os"
	"sync"
	"time"

	"github.com/gdnd-community/core-kit/pkg/meta"
	"github.com/rs/zerolog"
)

var (
	log      zerolog.Logger
	initOnce sync.Once
)

type Option func(*zerolog.Logger)

func mergeFields(fields ...map[string]any) map[string]interface{} {
	merged := make(map[string]any)
	for _, f := range fields {
		for k, v := range f {
			merged[k] = v
		}
	}
	return merged
}

func WithMetadata(meta *meta.Metadata) Option {
	return func(l *zerolog.Logger) {
		*l = l.With().
			Str("hostname", meta.Hostname).
			Str("app_name", meta.AppName).
			Str("app_version", meta.AppVersion).
			Str("env", meta.Env).Logger()
	}
}

func WithDevelopmentMode() Option {
	return func(l *zerolog.Logger) {
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		*l = l.Output(consoleWriter)
	}
}

func Init(level string, opts ...Option) {
	initOnce.Do(func() {
		zerolog.TimeFieldFormat = time.RFC3339

		lvl, err := zerolog.ParseLevel(level)
		if err != nil {
			lvl = zerolog.InfoLevel
		}

		// Log seviyesi Debug veya altındaysa caller ekler.
		var loggerContext zerolog.Context
		if lvl <= zerolog.DebugLevel {
			loggerContext = zerolog.New(os.Stdout).Level(lvl).With().Timestamp().Caller()
		} else {
			loggerContext = zerolog.New(os.Stdout).Level(lvl).With().Timestamp()
		}

		baseLogger := loggerContext.Logger()

		for _, opt := range opts {
			opt(&baseLogger)
		}

		log = baseLogger
	})
}

func Info(msg string, fields ...map[string]any) {
	if len(fields) > 0 {
		mergedFields := mergeFields(fields...)
		log.Info().Fields(mergedFields).Msg(msg)
	} else {
		log.Info().Msg(msg)
	}
}

func Warn(msg string, fields ...map[string]any) {
	if len(fields) > 0 {
		mergedFields := mergeFields(fields...)
		log.Warn().Fields(mergedFields).Msg(msg)
	} else {
		log.Warn().Msg(msg)
	}
}

func Debug(msg string, fields ...map[string]any) {
	if len(fields) > 0 {
		mergedFields := mergeFields(fields...)
		log.Debug().Fields(mergedFields).Msg(msg)
	} else {
		log.Debug().Msg(msg)
	}
}

func Error(err error, msg string, fields ...map[string]any) {
	event := log.Error().Err(err)
	if len(fields) > 0 {
		mergedFields := mergeFields(fields...)
		event.Fields(mergedFields).Msg(msg)
	} else {
		event.Msg(msg)
	}
}
