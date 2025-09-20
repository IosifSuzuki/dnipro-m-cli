package logger

import (
	"errors"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Level int8

const (
	// LevelDebug logs are typically voluminous, and are usually disabled in
	// production.
	LevelDebug Level = iota - 1
	// LevelInfo is the default logging priority.
	LevelInfo
	// LevelWarn logs are more important than Info, but don't need individual
	// human review.
	LevelWarn
	// LevelError logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	LevelError
	// LevelPanic logs a message, then panics.
	LevelPanic
	// LevelFatal logs a message, then calls os.Exit(1).
	LevelFatal
)

func (l *Level) FromString(lvl string) error {
	switch strings.TrimSpace(strings.ToLower(lvl)) {
	case "debug":
		*l = LevelDebug
	case "info":
		*l = LevelInfo
	case "warn":
		*l = LevelWarn
	case "error":
		*l = LevelError
	case "panic":
		*l = LevelPanic
	case "fatal":
		*l = LevelFatal
	default:
		return errors.New("invalid log level")
	}
	return nil
}

type ENV int8

const (
	DEV ENV = iota
	PROD
)

func ENVFromString(env string) (ENV, error) {
	var e ENV
	switch strings.TrimSpace(strings.ToLower(env)) {
	case "prod":
		e = PROD
	case "dev":
		e = DEV
	default:
		return DEV, errors.New("invalid log environment")
	}
	return e, nil
}

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Panic(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	Log(lvl Level, msg string, fields ...Field)
	Flush() error
}

type Field struct {
	Key string
	Val interface{}
}

func F(key string, val interface{}) Field {
	return Field{Key: key, Val: val}
}

func FError(e error) Field {
	return Field{Key: "error", Val: e}
}

type logger struct {
	lg *zap.Logger
}

func (l logger) Debug(msg string, fields ...Field) {
	l.Log(LevelDebug, msg, fields...)
}

func (l logger) Info(msg string, fields ...Field) {
	l.Log(LevelInfo, msg, fields...)
}

func (l logger) Warn(msg string, fields ...Field) {
	l.Log(LevelWarn, msg, fields...)
}

func (l logger) Error(msg string, fields ...Field) {
	l.Log(LevelError, msg, fields...)
}

func (l logger) Panic(msg string, fields ...Field) {
	l.Log(LevelPanic, msg, fields...)
}

func (l logger) Fatal(msg string, fields ...Field) {
	l.Log(LevelFatal, msg, fields...)
}

func (l logger) Flush() error {
	return l.lg.Sync()
}

func (l logger) Log(lvl Level, msg string, fields ...Field) {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key, field.Val)
	}
	l.lg.Log(zapcore.Level(lvl), msg, zapFields...)
}

func NewLogger(env ENV) Logger {
	var level Level
	switch env {
	case PROD:
		level = LevelInfo
	default:
		level = LevelDebug
	}
	return logger{
		lg: buildLogger(env, zapcore.Level(level)),
	}
}

func buildLogger(env ENV, level zapcore.Level) *zap.Logger {
	stdout := zapcore.AddSync(os.Stdout)
	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
	})

	switch env {
	case PROD:
		productionCfg := zap.NewProductionEncoderConfig()
		productionCfg.TimeKey = "timestamp"
		productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

		consoleEncoder := zapcore.NewConsoleEncoder(productionCfg)
		fileEncoder := zapcore.NewJSONEncoder(productionCfg)
		core := zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, stdout, level),
			zapcore.NewCore(fileEncoder, file, level),
		)
		return zap.New(core)
	case DEV:
		developmentCfg := zap.NewDevelopmentEncoderConfig()
		developmentCfg.TimeKey = "timestamp"
		developmentCfg.EncodeTime = zapcore.ISO8601TimeEncoder

		consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
		fileEncoder := zapcore.NewConsoleEncoder(developmentCfg)

		core := zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, stdout, level),
			zapcore.NewCore(fileEncoder, file, level),
		)
		return zap.New(core)
	}
	return nil
}
