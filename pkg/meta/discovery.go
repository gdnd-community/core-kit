package meta

import (
	"crypto/rand"
	"os"
	"runtime"
	"time"
)

type Metadata struct {
	Hostname   string        `json:"hostname"`
	PodName    string        `json:"pod_name"`
	Namespace  string        `json:"namespace"`
	NodeName   string        `json:"node_name"`
	AppName    string        `json:"app_name"`
	AppVersion string        `json:"app_version"`
	GoVersion  string        `json:"go_version"`
	Env        string        `json:"env"`
	StartTime  string        `json:"start_time"`
	Uptime     time.Duration `json:"uptime"`
	InstanceID string        `json:"instance_id"`
}

var startTime = time.Now().Format("2006-01-02 15:04:05")

func Discover(appName, appVersion, env string) *Metadata {
	hostname, _ := os.Hostname()
	layout := "2006-01-02 15:04:05"
	since, err := time.Parse(layout, startTime)
	if err != nil {
		return nil
	}

	podName := getenv("POD_NAME", "")
	var instanceID string
	if podName != "" {
		instanceID = podName
	} else {
		instanceID = shortID(8)
	}

	return &Metadata{
		Hostname:   hostname,
		PodName:    podName,
		Namespace:  getenv("POD_NAMESPACE", "default"),
		NodeName:   getenv("NODE_NAME", "unknown"),
		AppName:    appName,
		AppVersion: appVersion,
		GoVersion:  runtime.Version(),
		Env:        env,
		StartTime:  startTime,
		Uptime:     time.Since(since),
		InstanceID: instanceID,
	}
}

func getenv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func shortID(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	rand.Read(b)
	for i := 0; i < n; i++ {
		b[i] = letters[int(b[i])%len(letters)]
	}
	return string(b)
}
