package main

import (
	// std

	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	// internal pkg

	"github.com/vishenosik/file-svc-sdk/api"
	"github.com/vishenosik/gocherry"
	_ctx "github.com/vishenosik/gocherry/pkg/context"
	"github.com/vishenosik/gocherry/pkg/errors"
	"github.com/vishenosik/gocherry/pkg/grpc"

	// internal

	"github.com/vishenosik/file-svc/internal/dto"
	"github.com/vishenosik/file-svc/internal/store/config"
	"github.com/vishenosik/file-svc/internal/store/mongodb"
	"github.com/vishenosik/file-svc/internal/usecase"
)

func main() {

	gocherry.Flags(os.Stdout, os.Args[1:],
		gocherry.AppFlags(os.Stdout),
		gocherry.ConfigFlags(os.Stdout),
	)

	flag.Parse()
	ctx := context.Background()

	app, err := NewApp()
	if err != nil {
		log.Fatalf("failed to init app %s", err.Error())
	}

	err = app.Start(ctx)
	if err != nil {
		log.Fatalf("failed to start app %s", err.Error())
	}

	// Graceful shut down
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	stopctx, cancel := context.WithTimeout(
		_ctx.WithStopCtx(ctx, <-stop),
		time.Second*5,
	)
	defer cancel()

	app.Stop(stopctx)
}

func NewApp() (*gocherry.App, error) {

	// STORES

	mongoStore, err := mongodb.NewFileStore()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create file store")
	}

	configStore, err := config.NewService()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create config store")
	}

	// localStore := local.NewFileStore()

	// USECASES

	settingsUsc := usecase.NewService(configStore)

	providerUsc := usecase.NewProvider(mongoStore)

	infoUsc := dto.NewInfoDTO(usecase.NewInfo(mongoStore))

	// API

	fileService := api.NewFileServiceApi(providerUsc, infoUsc, settingsUsc)

	// SERVICES

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
		mongoStore,
	)

	return app, nil
}
