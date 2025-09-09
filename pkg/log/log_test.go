package log

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"regexp"
	"sync"
	"testing"

	"github.com/gdnd-community/core-kit/pkg/meta"
	"github.com/stretchr/testify/assert"
)

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stdout = old
	return buf.String()
}

func Test_Logger(t *testing.T) {
	metas := meta.Discover("CoreKit Test", "1.0.0", "Local")

	t.Run("Info log should be JSON and contain metadata", func(t *testing.T) {
		initOnce = sync.Once{}
		output := captureOutput(func() {
			Init("info", WithMetadata(metas))
			Info("Info test message", map[string]any{"key1": "value1"})
		})

		var logEntry map[string]interface{}
		err := json.Unmarshal([]byte(output), &logEntry)
		assert.NoError(t, err, "Log output should be valid JSON")

		assert.Equal(t, "info", logEntry["level"], "Log level should be 'info'")
		assert.Equal(t, "Info test message", logEntry["message"], "Log message should be correct")
		assert.Equal(t, "CoreKit Test", logEntry["app_name"], "Metadata should be passed correctly")
		assert.Equal(t, "value1", logEntry["key1"], "Custom fields should be added correctly")
	})

	t.Run("Debug log in JSON format should contain caller info", func(t *testing.T) {
		initOnce = sync.Once{}
		output := captureOutput(func() {
			Init("debug")
			Debug("Debug test message")
		})

		var logEntry map[string]interface{}
		err := json.Unmarshal([]byte(output), &logEntry)
		assert.NoError(t, err, "Debug log should be valid JSON")

		assert.Equal(t, "debug", logEntry["level"], "Log level should be 'debug'")
		assert.Contains(t, logEntry, "caller", "Debug log should include caller information")

		caller, ok := logEntry["caller"].(string)
		assert.True(t, ok, "Caller field should be a string")
		assert.Regexp(t, regexp.MustCompile(`.*log\.go:\d+`), caller, "Caller format should be 'any/path/log.go:line'")
	})

	t.Run("Error log should contain error field", func(t *testing.T) {
		initOnce = sync.Once{}
		output := captureOutput(func() {
			Init("info")
			err := errors.New("sample error")
			Error(err, "An error occurred", map[string]any{"code": 500})
		})

		var logEntry map[string]interface{}
		err := json.Unmarshal([]byte(output), &logEntry)
		assert.NoError(t, err, "Error log should be valid JSON")

		assert.Equal(t, "error", logEntry["level"], "Log level should be 'error'")
		assert.Equal(t, "sample error", logEntry["error"], "Error field should be correct")
		assert.Equal(t, "An error occurred", logEntry["message"], "Message field should be correct")
	})

	t.Run("Development mode output should be human-readable with caller info", func(t *testing.T) {
		initOnce = sync.Once{}
		output := captureOutput(func() {
			Init("debug", WithDevelopmentMode())
			Info("Info log in dev mode")
			Debug("Debug log in dev mode")
		})

		var logEntry map[string]interface{}
		err := json.Unmarshal([]byte(output), &logEntry)
		assert.Error(t, err, "Development mode output should not be JSON")

		assert.Contains(t, output, "INF", "Info log should have 'INF' tag")
		assert.Contains(t, output, "DBG", "Debug log should have 'DBG' tag")
		assert.Contains(t, output, "Info log in dev mode", "Info message should be present")
		assert.Contains(t, output, "Debug log in dev mode", "Debug message should be present")

		assert.Regexp(t, regexp.MustCompile(`DBG.*log\.go:\d+`), output, "Debug log should include caller info in human-readable format")
	})

	t.Run("Logs below level should not be printed", func(t *testing.T) {
		initOnce = sync.Once{}
		output := captureOutput(func() {
			Init("warn")
			Info("This info message should not appear")
			Debug("This debug message should not appear")
		})

		assert.Empty(t, output, "Logs below the configured level should not be printed")
	})
}

func Benchmark_Logger_Simple(b *testing.B) {
	initOnce = sync.Once{}
	Init("info")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info("hello")
	}
}

func Benchmark_Logger_WithFields(b *testing.B) {
	Init("debug")

	fields := map[string]any{
		"request_id":  "12345",
		"endpoint":    "/api/v1/users",
		"duration_ms": 150,
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Debug("some message", fields)
	}
}
