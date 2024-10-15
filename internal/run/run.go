package run

import (
	"bmp-tgbot/internal/sdk"
	"go.uber.org/zap"
	"os"
)

var Logger *zap.Logger

func Init() {
	cfg := zap.NewProductionConfig()
	if debug := os.Getenv(sdk.EnvDebug); debug != "" {
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	Logger = logger
}
