package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var (
	Log *logger
)

const (
	envDevelopment = "dev"
	envProduction  = "prod"
)

type logger struct {
	*zap.SugaredLogger
}

func New(debug bool, timeZone string, envType string, lokiConfig *LokiConfig) {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "timestamp",
		CallerKey:      "caller",
		EncodeTime:     customTimeEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	if timeZone != "" {
		encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.In(time.FixedZone(timeZone, 3*60*60)).Format("2006-01-02 15:04:05"))
		}
	}

	var level zapcore.Level
	if debug {
		level = zapcore.DebugLevel
	} else {
		level = zapcore.InfoLevel
	}

	var encoder zapcore.Encoder
	if envType == envProduction {
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	sinks := []zapcore.WriteSyncer{zapcore.Lock(os.Stdout)} // stdout

	if lokiConfig != nil && lokiConfig.url != "" {
		sinks = append(sinks, newLokiSink(*lokiConfig))
	}

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(sinks...), level),
	)
	log := zap.New(core, zap.AddCaller())

	Log = &logger{
		SugaredLogger: log.Sugar(),
	}
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.In(time.FixedZone("GMT+0", 3*60*60)).Format("2006-01-02 15:04:05"))
}
