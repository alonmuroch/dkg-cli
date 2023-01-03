package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var Logger = func() *zap.Logger {
	cfg := zap.Config{
		Encoding:    "console",
		Level:       zap.NewAtomicLevelAt(zapcore.InfoLevel),
		OutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",
			//LevelKey:    "level",
			EncodeLevel: zapcore.CapitalColorLevelEncoder,
			//TimeKey:     "time",
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.UTC().Format("15:04:05")) // iso3339CleanTime
			},
			CallerKey:      "caller",
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		},
	}

	ret, err := cfg.Build()
	if err != nil {
		panic(err.Error())
	}

	return ret
}
