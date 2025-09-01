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

func Info(msg string, fields map[string]interface{}) {
	log.Info().Fields(fields).Msg(msg)
}

func Error(msg string, fields map[string]interface{}) {
	log.Error().Fields(fields).Msg(msg)
}

func Debug(msg string, fields map[string]interface{}) {
	log.Debug().Fields(fields).Msg(msg)
}

func Warn(msg string, fields map[string]interface{}) {
	log.Warn().Fields(fields).Msg(msg)
}
