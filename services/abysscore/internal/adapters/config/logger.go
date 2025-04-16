package config

import "github.com/intezya/pkglib/logger"

type LoggerConfig struct {
	lokiConfig *logger.LokiConfig
}

func SetupLogger(isDebug bool, env string, config *LoggerConfig) {
	_, err := logger.New(
		logger.WithDebug(isDebug),
		logger.WithCaller(true),
		logger.WithLoki(config.lokiConfig),
		logger.WithEnvironment(env),
	)

	logger.Log.Debugf("Debug mode: %t", isDebug)

	if err != nil {
		panic(err)
	}
}
