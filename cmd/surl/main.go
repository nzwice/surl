package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nzwice/surl/pkg/config"
	"github.com/nzwice/surl/pkg/endpoints"
	"github.com/nzwice/surl/pkg/logging"
	logger "github.com/nzwice/surl/pkg/logging"
	"github.com/nzwice/surl/pkg/shorten"
	"github.com/nzwice/surl/pkg/surldb"
	"github.com/nzwice/surl/pkg/transport"
	"github.com/uptrace/bun"
	"golang.org/x/sync/errgroup"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "")

	ctx := context.Background()
	ctx, cancelFunc := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancelFunc()

	logger.SetupLogger()

	cfg, err := config.Load(configPath)
	if err != nil {
		slog.ErrorContext(ctx, "fail to load config", logger.ErrorAttr(err))
		return
	}

	var db *bun.DB
	{
		db, err = surldb.New(cfg.DB, cfg.Debug)
		if err != nil {
			slog.ErrorContext(ctx, "fail to connect to db", logger.ErrorAttr(err))
			return
		}
	}

	var shortenSvc shorten.Service
	{
		shortenSvc = shorten.New(db)
	}

	endpoints := endpoints.MakeEndpoints(shortenSvc)

	httpHandler := transport.HttpHandler(endpoints)
	httpServer := &http.Server{
		Addr:    cfg.HttpAddr,
		Handler: httpHandler,
	}

	var wg = new(errgroup.Group)

	wg.Go(func() error {
		slog.InfoContext(ctx, "starting http server...")
		return httpServer.ListenAndServe()
	})

	wg.Go(func() error {
		<-ctx.Done()
		slog.InfoContext(ctx, "shutting down http server...")

		shutdownCtx, cancelFunc := context.WithTimeout(ctx, 5*time.Second)
		defer cancelFunc()

		return httpServer.Shutdown(shutdownCtx)
	})

	if err := wg.Wait(); err != nil {
		slog.ErrorContext(ctx, "server exit", logging.ErrorAttr(err))
	}
}
