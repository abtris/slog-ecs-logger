package main_test

import (
	"testing"

	logger "github.com/abtris/slog-ecs-logger"
)

func TestJsonHandler(t *testing.T) {
	log := logger.GetLogger()
	log.Info("test")
}
