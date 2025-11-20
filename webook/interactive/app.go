package main

import (
	"Book_Exp/webook/pkg/grpcx"
	"Book_Exp/webook/pkg/saramax"
)

type App struct {
	server    *grpcx.Server
	consumers []saramax.Consumer
}
