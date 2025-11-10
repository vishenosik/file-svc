package main

import (
	// std

	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	// internal pkg
	"github.com/vishenosik/file-svc-sdk/api"
	"github.com/vishenosik/file-svc/internal/store"
	"github.com/vishenosik/file-svc/internal/usecases"
	"github.com/vishenosik/gocherry"
	_ctx "github.com/vishenosik/gocherry/pkg/context"
	"github.com/vishenosik/gocherry/pkg/grpc"
	// internal
)

func main() {

	gocherry.Flags(os.Stdout, os.Args[1:],
		gocherry.AppFlags(os.Stdout),
		gocherry.ConfigFlags(os.Stdout),
	)

	flag.Parse()
	ctx := context.Background()

	// App init
	app, err := NewApp()
	if err != nil {
		panic(err)
	}

	err = app.Start(ctx)
	if err != nil {
		panic(err)
	}

	// Graceful shut down
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	stopctx, cancel := context.WithTimeout(
		_ctx.WithStopCtx(ctx, <-stop),
		time.Second*5,
	)
	defer cancel()

	if err := app.Stop(stopctx); err != nil {
		panic(err)
	}
}

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type App struct {
	Server
}

func NewApp() (*App, error) {

	// STORES

	fileStore, err := store.NewFileStore()
	if err != nil {
		return nil, err
	}

	// USECASES

	fileUsecase, err := usecases.NewFileService(
		fileStore,
		fileStore,
	)
	if err != nil {
		return nil, err
	}

	// API

	fileService := api.NewFileServiceApi(fileUsecase)

	grpcServer, err := grpc.NewGrpcServer(
		grpc.GrpcServices{
			fileService,
		},
		grpc.WithLogInterceptors(),
	)
	if err != nil {
		return nil, err
	}

	app, err := gocherry.NewApp()
	if err != nil {
		return nil, err
	}

	app.AddServices(
		grpcServer,
	)

	return &App{
		Server: app,
	}, nil
}
