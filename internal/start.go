package internal

import (
	"context"
	"github.com/coffeehc/xnetforward/internal/services/forwardservice"

	"github.com/coffeehc/boot/engine"
	"github.com/spf13/cobra"
)

func Start(ctx context.Context, cmd *cobra.Command, args []string) (engine.ServiceCloseCallback, error) {
	forwardservice.EnablePlugin(ctx)
	return func() {

	}, nil
}
