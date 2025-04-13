package main

import (
	"context"
	"flag"
	"os/signal"
	"sync"
	"syscall"

	"github.com/mohammadne/gorillamq/internal/config"
	"github.com/mohammadne/gorillamq/internal/manager"
	"github.com/mohammadne/gorillamq/pkg/logger"
	"github.com/mohammadne/gorillamq/pkg/tcp"
)

func main() {
	monitorPort := flag.Int("monitor-port", 8001, "The server port which handles monitoring endpoints (default: 8001)")
	flag.Parse() // Parse the command-line flags

	cfg := config.Load(true)
	logger := logger.NewZap(cfg.Logger)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	var wg sync.WaitGroup

	wg.Add(1)
	tcp := tcp.NewTCP(cfg.TCP)
	manager.NewBroker(logger, tcp).Start(ctx)
	_ = monitorPort

	<-ctx.Done()
	wg.Wait()
	logger.Warn("interruption signal recieved, gracefully shutdown the server")
}

// import (
// 	"context"
// 	"os/signal"
// 	"syscall"

// 	"github.com/mohammadne/gorillamq/internal/config"
// 	"github.com/mohammadne/gorillamq/internal/manager"
// 	"github.com/mohammadne/gorillamq/pkg/logger"
// 	"github.com/mohammadne/gorillamq/pkg/tcp"
// 	"github.com/spf13/cobra"
// 	"go.uber.org/zap"
// )

// type Server struct {
// 	config *config.Config
// 	logger *zap.Logger
// }

// func (server Server) Command() *cobra.Command {
// 	run := func(_ *cobra.Command, _ []string) {
// 		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
// 		defer stop()

// 		server.config = config.Load(true)
// 		server.logger = logger.NewZap(server.config.Logger)

// 		server.start(ctx)

// 		<-ctx.Done()
// 		server.logger.Warn("Got interruption signal, gracefully shutdown the server")

// 		server.stop(ctx)
// 	}

// 	cmd := &cobra.Command{
// 		Use:   "server",
// 		Short: "Run GorillaMQ server",
// 		Run:   run,
// 	}

// 	// cmd.Flags().IntVar(&server.ports.management, "management-port", 8000, "The port the metric and probe endpoints binds to")

// 	return cmd
// }

// func (server *Server) start(ctx context.Context) {
// 	tcp := tcp.NewTCP(server.config.TCP)
// 	manager.NewBroker(server.logger, tcp).Start(ctx)
// }

// func (server *Server) stop(ctx context.Context) {}
