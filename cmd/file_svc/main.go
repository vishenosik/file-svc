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
	"github.com/vishenosik/gocherry"
	_ctx "github.com/vishenosik/gocherry/pkg/context"
	"github.com/vishenosik/gocherry/pkg/errors"
	"github.com/vishenosik/gocherry/pkg/grpc"
	"github.com/vishenosik/gocherry/pkg/logs"

	// internal

	"github.com/vishenosik/file-svc/internal/dto"
	"github.com/vishenosik/file-svc/internal/store/config"
	"github.com/vishenosik/file-svc/internal/store/mongodb"
	"github.com/vishenosik/file-svc/internal/usecase"
)

func main() {

	log := logs.SetupLogger().With(logs.AppComponent("main"))

	gocherry.Flags(os.Stdout, os.Args[1:],
		gocherry.AppFlags(os.Stdout),
		gocherry.ConfigFlags(os.Stdout),
	)

	flag.Parse()

	ctx := context.Background()

	app, err := NewApp(ctx)
	if err != nil {
		log.Error("failed to init app", logs.Error(err))
		os.Exit(1)
	}

	err = app.Start(ctx)
	if err != nil {
		log.Error("failed to start app", logs.Error(err))
		os.Exit(1)
	}

	// Graceful shut down
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	stopctx, stopCancel := context.WithTimeout(_ctx.WithStopCtx(context.Background(), <-stop), time.Second*5)
	defer stopCancel()

	app.Stop(stopctx)
}

func NewApp(ctx context.Context) (*gocherry.App, error) {

	// return nil, errors.New("testing")

	// STORES

	mongoStore, err := mongodb.NewFileStoreRetry()
	if err != nil {
		return nil, err
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
