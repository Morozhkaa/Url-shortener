package logger

import (
	"fmt"

	"go.uber.org/zap"
)

type Logger = *zap.Logger

type LoggerOptions struct {
	IsProd bool
}

func New(opts LoggerOptions) (Logger, error) {
	var zapCfg zap.Config

	if opts.IsProd {
		zapCfg = zap.NewProductionConfig()
	} else {
		zapCfg = zap.NewDevelopmentConfig()
	}

	l, err := zapCfg.Build()
	if err != nil {
		return nil, fmt.Errorf("create logger failed: %w", err)
	}
	zap.ReplaceGlobals(l)
	return l, nil
}
