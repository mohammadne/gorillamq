package cmd

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/mohammadne/gorillamq/internal/config"
	"github.com/mohammadne/gorillamq/internal/manager"
	"github.com/mohammadne/gorillamq/pkg/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type Server struct {
	metricsPort int
	config      *config.Config
	logger      *zap.Logger
}

func (server Server) Command() *cobra.Command {
	run := func(_ *cobra.Command, _ []string) {
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		server.config = config.Load(true)
		server.logger = logger.NewZap(server.config.Logger)

		server.start(ctx)

		<-ctx.Done()
		server.logger.Warn("Got interruption signal, gracefully shutdown the server")

		server.stop(ctx)
	}

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run GorillaMQ server",
		Run:   run,
	}

	cmd.Flags().IntVar(&server.metricsPort, "metrics-port", 8081, "The port the metric and probe endpoints binds to")

	return cmd
}

func (server *Server) start(ctx context.Context) {
	manager.NewBroker(server.logger).Start(ctx, server.config.Network)
}

func (server *Server) stop(ctx context.Context) {}
