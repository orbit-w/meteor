package transport

import (
	mlog "github.com/orbit-w/meteor/modules/mlog_v2"
)

var (
	logger = mlog.NewFileLogger(
		mlog.WithDevelopment(true),
		mlog.WithOutputPaths("logs/transport.log"),
		mlog.WithRotation(100, 3, 3, false),
		mlog.WithFormat("console"))
)

func SetLogger(log *mlog.Logger) {
	logger = log
}
