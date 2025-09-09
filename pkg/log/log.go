package log

import (
	"fmt"
	"os"
	"runtime"
	"strings"
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

func mergeFields(fields ...map[string]any) map[string]any {
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

		devMode := false
		for _, opt := range opts {
			if fmt.Sprintf("%T", opt) == "log.WithDevelopmentMode" {
				devMode = true
				break
			}
		}

		if devMode {
			zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
				pcs := make([]uintptr, 10)
				n := runtime.Callers(4, pcs)
				frames := runtime.CallersFrames(pcs[:n])

				for {
					f, more := frames.Next()
					if !strings.Contains(f.File, "github.com/gdnd-community/core-kit/pkg/log/") &&
						!strings.Contains(f.File, "github.com/rs/zerolog") {
						short := f.File
						if idx := strings.LastIndex(f.File, "/"); idx != -1 {
							short = f.File[idx+1:]
						}
						return fmt.Sprintf("%s:%d", short, f.Line)
					}
					if !more {
						break
					}
				}
				return fmt.Sprintf("%s:%d", file, line)
			}
		}

		lvl, err := zerolog.ParseLevel(level)
		if err != nil {
			lvl = zerolog.InfoLevel
		}

		var loggerContext zerolog.Context
		if devMode {
			loggerContext = zerolog.New(os.Stdout).
				Level(lvl).
				With().
				Timestamp().
				CallerWithSkipFrameCount(0)
		} else if lvl <= zerolog.DebugLevel {
			loggerContext = zerolog.New(os.Stdout).
				Level(lvl).
				With().
				Timestamp()
		} else {
			loggerContext = zerolog.New(os.Stdout).
				Level(lvl).
				With().
				Timestamp()
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
		log.Info().Fields(mergeFields(fields...)).Msg(msg)
	} else {
		log.Info().Msg(msg)
	}
}

func Warn(msg string, fields ...map[string]any) {
	if len(fields) > 0 {
		log.Warn().Fields(mergeFields(fields...)).Msg(msg)
	} else {
		log.Warn().Msg(msg)
	}
}

func Debug(msg string, fields ...map[string]any) {
	if len(fields) > 0 {
		log.Debug().Fields(mergeFields(fields...)).Msg(msg)
	} else {
		log.Debug().Msg(msg)
	}
}

func Error(err error, msg string, fields ...map[string]any) {
	event := log.Error().Err(err)
	if len(fields) > 0 {
		event.Fields(mergeFields(fields...)).Caller().Msg(msg)
	} else {
		event.Caller().Msg(msg)
	}
}
