package logger

import (
	"sync"

	"go.uber.org/zap"
)

var (
	log *zap.Logger
	mu  sync.Mutex
)

func Init(level, format string) error {
	mu.Lock()
	defer mu.Unlock()

	zapLevel := zap.InfoLevel
	switch level {
	case "debug":
		zapLevel = zap.DebugLevel
	case "warn", "warning":
		zapLevel = zap.WarnLevel
	case "error":
		zapLevel = zap.ErrorLevel
	}

	var encoding string
	if format == "json" {
		encoding = "json"
	} else {
		encoding = "console"
	}

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapLevel),
		Development:      false,
		Encoding:         encoding,
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	var err error
	log, err = config.Build()
	if err != nil {
		return err
	}

	return nil
}

func GetLogger() *zap.Logger {
	if log == nil {
		mu.Lock()
		defer mu.Unlock()
		if log == nil {
			// 默认配置
			config := zap.Config{
				Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
				Development:      false,
				Encoding:         "json",
				EncoderConfig:    zap.NewProductionEncoderConfig(),
				OutputPaths:      []string{"stdout"},
				ErrorOutputPaths: []string{"stderr"},
			}
			log, _ = config.Build()
		}
	}
	return log
}

func Debug(message string, fields map[string]interface{}) {
	GetLogger().Debug(message, toZapFields(fields)...)
}

func Info(message string, fields map[string]interface{}) {
	GetLogger().Info(message, toZapFields(fields)...)
}

func Warn(message string, fields map[string]interface{}) {
	GetLogger().Warn(message, toZapFields(fields)...)
}

func Error(message string, fields map[string]interface{}) {
	GetLogger().Error(message, toZapFields(fields)...)
}

func Fatal(message string, fields map[string]interface{}) {
	GetLogger().Fatal(message, toZapFields(fields)...)
}

func Sync() {
	if log != nil {
		log.Sync()
	}
}

func toZapFields(m map[string]interface{}) []zap.Field {
	fields := make([]zap.Field, 0, len(m))
	for k, v := range m {
		fields = append(fields, zap.Any(k, v))
	}
	return fields
}

// WithFields 返回一个新的logger，带有指定字段
func WithFields(fields map[string]interface{}) *zap.Logger {
	return GetLogger().With(toZapFields(fields)...)
}
