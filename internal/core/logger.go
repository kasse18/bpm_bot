package core

import (
	"fmt"

	"go.uber.org/zap"
)

type TelegramLogger struct {
	logger *zap.Logger
}

func (r *TelegramLogger) Println(v ...interface{}) {
	r.logger.Debug(fmt.Sprint(v...))
}

func (r *TelegramLogger) Printf(format string, v ...interface{}) {
	r.logger.Debug(fmt.Sprintf(format, v...))
}
