package mongodb

import (
	"context"
	"log/slog"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"github.com/vishenosik/gocherry/pkg/config"
	"github.com/vishenosik/gocherry/pkg/logs"
	"github.com/vishenosik/gocherry/pkg/retry"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const name = "photos"

var (
	logger  = logs.SetupLogger().With(logs.AppComponent("mongo"))
	conf    FileStoreConfig
	ErrConn = errors.New("failed to connect to db")
	ErrPing = errors.New("failed to ping db")
)

type FileStoreConfig struct {
	Uri      string `env:"MONGODB_URI"`
	Database string `env:"MONGODB_DATABASE"`
}

func (FileStoreConfig) Desc() string {
	return "File service config"
}

func (c FileStoreConfig) validate() error {
	_, err := url.Parse(c.Uri)
	if err != nil {
		return errors.Wrap(err, "failed to parse mongo uri")
	}
	return nil
}

type FileStore struct {
	client *mongo.Client
}

func NewFileStoreRetry() (*FileStore, error) {

	if err := loadConfig(); err != nil {
		return nil, err
	}

	uri, _ := url.Parse(conf.Uri)
	log := logger.With(
		slog.String("addr", uri.Host),
		slog.String("db", conf.Database),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var filestore *FileStore

	backoff := retry.NewFibonacci(1*time.Second, time.Minute)
	if err := retry.Do(ctx, backoff, func(ctx context.Context) error {

		log.Info("trying to connect to mongo")

		fs, err := NewFileStore(ctx)
		if err != nil {
			retryable := func() error {
				if errors.Is(err, context.DeadlineExceeded) {
					return retry.RetryableError(err)
				}
				if errors.Is(err, ErrConn) {
					return retry.RetryableError(err)
				}
				if errors.Is(err, ErrPing) {
					return retry.RetryableError(err)
				}
				return nil
			}

			msg := "failed to connect to mongo"

			if err := retryable(); err != nil {
				log.Error(msg,
					slog.Int64("retry_in_seconds", backoff.RetryInSeconds()),
					logs.Error(err),
				)
				return err
			}

			log.Error(msg, logs.Error(err))
			return err
		}

		log.Info("connected to mongo db")

		filestore = fs
		return nil

	}); err != nil {
		return nil, err
	}

	return filestore, nil
}

func (fs *FileStore) Close(ctx context.Context) error {
	return fs.client.Disconnect(ctx)
}

func NewFileStore(ctx context.Context) (*FileStore, error) {

	client, err := connect(ctx, conf.Uri)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to mongo")
	}

	return &FileStore{
		client: client,
	}, nil
}

func loadConfig() error {

	if err := config.ReadConfigEnv(&conf); err != nil {
		return errors.Wrap(err, "failed to read mongo config")
	}

	if err := conf.validate(); err != nil {
		return errors.Wrap(err, "failed to validate mongo config")
	}

	return nil
}

// To connect to mongodb
func connect(ctx context.Context, dsn string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// mongo.Connect return mongo.Client method
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn))
	if err != nil {
		return nil, ErrConn
	}

	// ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, ErrPing
	}

	return client, nil
}

func (fs *FileStore) bucket() (*gridfs.Bucket, error) {
	return gridfs.NewBucket(
		fs.client.Database(conf.Database),
		options.GridFSBucket().SetName(name),
	)
}
