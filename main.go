package main

import (
	"context"
	"github.com/coffeehc/boot/configuration"
	"github.com/coffeehc/boot/engine"
	"github.com/coffeehc/xnetforward/internal"
)

func main() {
	if configuration.Version == "" {
		configuration.Version = "0.0.1"
	}
	engine.StartEngine(context.TODO(), internal.GetServiceInfo(), internal.Start)
}
