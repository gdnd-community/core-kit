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

func mergeFields(fields ...map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})
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
			Str("pod_name", meta.PodName).
			Str("namespace", meta.Namespace).
			Str("node_name", meta.NodeName).
			Str("app_name", meta.AppName).
			Str("app_version", meta.AppVersion).
			Str("env", meta.Env).
			Str("instance_id", meta.InstanceID).
			Logger()
	}
}
func Init(level string, opts ...Option) {
	initOnce.Do(func() {
		zerolog.TimeFieldFormat = time.RFC3339

		lvl, err := zerolog.ParseLevel(level)
		if err != nil {
			lvl = zerolog.InfoLevel
		}

		baseLogger := zerolog.New(os.Stdout).
			Level(lvl).
			With().
			Timestamp().
			Logger()

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
