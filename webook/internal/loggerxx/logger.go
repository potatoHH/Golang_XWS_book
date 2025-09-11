package loggerxx

import "go.uber.org/zap"

var Logger *zap.Logger

func InitLogger() {
	var err error
	Logger, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
}
