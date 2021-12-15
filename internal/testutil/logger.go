package testutil

import (
	"go.uber.org/zap"
	"gotest.tools/v3/assert"
	"testing"
)

func NewLogger(t *testing.T) *zap.Logger {
	cfg := zap.NewDevelopmentConfig()
	cfg.Level.SetLevel(zap.FatalLevel)
	logger, err := cfg.Build()
	assert.NilError(t, err, "create logger")
	return logger
}
