package logger

import (
	"context"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.SugaredLogger
}

type ctxKey struct{}

var (
	globalLogger *Logger
)

func Init(level, format, output string) (*Logger, error) {
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder
	if format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	var writeSyncer zapcore.WriteSyncer
	if output == "stdout" || output == "" {
		writeSyncer = zapcore.AddSync(os.Stdout)
	} else {
		file, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		writeSyncer = zapcore.AddSync(file)
	}

	core := zapcore.NewCore(encoder, writeSyncer, zapLevel)
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	globalLogger = &Logger{zapLogger.Sugar()}
	return globalLogger, nil
}

func GetLogger() *Logger {
	if globalLogger == nil {
		globalLogger, _ = Init("info", "json", "stdout")
	}
	return globalLogger
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	traceID, _ := ctx.Value(ctxKey{}).(string)
	if traceID != "" {
		return &Logger{l.With("trace_id", traceID)}
	}
	return l
}

func (l *Logger) WithFields(fields ...interface{}) *Logger {
	return &Logger{l.With(fields...)}
}

func (l *Logger) Sync() error {
	return l.SugaredLogger.Sync()
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, ctxKey{}, traceID)
}

func GetTraceID(ctx context.Context) string {
	traceID, _ := ctx.Value(ctxKey{}).(string)
	return traceID
}

func Info(ctx context.Context, msg string, fields ...interface{}) {
	GetLogger().WithContext(ctx).Infow(msg, fields...)
}

func Debug(ctx context.Context, msg string, fields ...interface{}) {
	GetLogger().WithContext(ctx).Debugw(msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...interface{}) {
	GetLogger().WithContext(ctx).Warnw(msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...interface{}) {
	GetLogger().WithContext(ctx).Errorw(msg, fields...)
}

func Fatal(ctx context.Context, msg string, fields ...interface{}) {
	GetLogger().WithContext(ctx).Fatalw(msg, fields...)
}

func LogRequest(ctx context.Context, method, path string, statusCode int, duration time.Duration, err error) {
	fields := []interface{}{
		"method", method,
		"path", path,
		"status_code", statusCode,
		"duration_ms", duration.Milliseconds(),
	}
	
	if err != nil {
		fields = append(fields, "error", err.Error())
		Error(ctx, "HTTP request", fields...)
	} else {
		Info(ctx, "HTTP request", fields...)
	}
}
