package meta

import (
	"os"
	"time"
)

type Metadata struct {
	Hostname   string        `json:"hostname"`
	AppName    string        `json:"app_name"`
	AppVersion string        `json:"app_version"`
	Env        string        `json:"env"`
	Uptime     time.Duration `json:"uptime"`
}

func Discover(appName, appVersion, env string) *Metadata {
	hostname, _ := os.Hostname()

	return &Metadata{
		Hostname:   hostname,
		AppName:    appName,
		AppVersion: appVersion,
		Env:        env,
	}
}
