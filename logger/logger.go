package logger

import (
	"fmt"
	zp "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.ytsaurus.tech/library/go/core/log"
	"go.ytsaurus.tech/library/go/core/log/zap"
	"os"
)

func parseLevel(level string) levels {
	zpLvl := zapcore.InfoLevel
	lvl := log.InfoLevel
	if level != "" {
		fmt.Printf("overriden YT log level to: %v\n", level)
		var l zapcore.Level
		if err := l.UnmarshalText([]byte(level)); err == nil {
			zpLvl = l
		}
		var gl log.Level
		if err := gl.UnmarshalText([]byte(level)); err == nil {
			lvl = gl
		}
	}
	return levels{zpLvl, lvl}
}

func LogLevel() string {
	if os.Getenv("LOG_LEVEL") != "" {
		return os.Getenv("LOG_LEVEL")
	}
	return "INFO"
}

type levels struct {
	Zap zapcore.Level
	Log log.Level
}

func getEnvLogLevels() levels {
	if level, ok := os.LookupEnv("LOG_LEVEL"); ok {
		return parseLevel(level)
	}
	return levels{zapcore.InfoLevel, log.InfoLevel}
}

func NewFileLog(logFile string) log.Logger {
	consoleLevel := getEnvLogLevels()
	defaultPriority := levelEnablerFactory(consoleLevel.Zap)
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	fileSync := zapcore.AddSync(file)
	stdErrEncoder := zapcore.NewConsoleEncoder(zap.CLIConfig(consoleLevel.Log).EncoderConfig)
	lbCore := zapcore.NewTee(
		zapcore.NewCore(stdErrEncoder, fileSync, defaultPriority),
	)

	return newLogger(lbCore)
}

func NewConsoleLogger() log.Logger {
	consoleLevel := getEnvLogLevels()
	defaultPriority := levelEnablerFactory(consoleLevel.Zap)
	syncStderr := zapcore.AddSync(os.Stderr)
	stdErrEncoder := zapcore.NewConsoleEncoder(zap.CLIConfig(consoleLevel.Log).EncoderConfig)
	lbCore := zapcore.NewTee(
		zapcore.NewCore(stdErrEncoder, syncStderr, defaultPriority),
	)

	return newLogger(lbCore)
}

func newLogger(core zapcore.Core) log.Logger {
	return &zap.Logger{
		L: zp.New(
			core,
			zp.AddCaller(),
			zp.AddCallerSkip(1),
			zp.AddStacktrace(zp.WarnLevel),
		),
	}
}

func levelEnablerFactory(zapLvl zapcore.Level) zapcore.LevelEnabler {
	return zp.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= zapLvl
	})
}
