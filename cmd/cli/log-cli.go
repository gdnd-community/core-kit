package main

import "github.com/gdnd-community/core-kit/pkg/log"

func main() {
	log.Init("info", log.WithDevelopmentMode())
	log.Info("system test", map[string]any{
		"data": "xxx",
	})
}
