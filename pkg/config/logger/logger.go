package logger

import (
	"log/slog"
	"os"

	"github.com/OzkrOssa/ros-iface-streamer/pkg/config"
	slogmulti "github.com/samber/slog-multi"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Set(config *config.App) {
	var logger *slog.Logger
	logger = slog.New(
		slog.NewTextHandler(os.Stderr, nil),
	)

	if config.Env == "prod" {
		logRotate := &lumberjack.Logger{
			Filename:   "log/app.log",
			MaxSize:    100, // megabytes
			MaxBackups: 3,
			MaxAge:     28, // days
			Compress:   true,
		}

		logger = slog.New(
			slogmulti.Fanout(
				slog.NewJSONHandler(logRotate, nil),
				slog.NewTextHandler(os.Stderr, nil),
			),
		)
	}

	slog.SetDefault(logger)
}
