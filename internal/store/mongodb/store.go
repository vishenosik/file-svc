package mongodb

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/vishenosik/gocherry/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const name = "photos"

type FileStoreConfig struct {
	Uri      string `env:"MONGODB_URI"`
	Database string `env:"MONGODB_DATABASE"`
}

func (FileStoreConfig) Desc() string {
	return "File service config"
}

func (c FileStoreConfig) validate() error {
	return nil
}

type FileStore struct {
	client *mongo.Client
	config FileStoreConfig
}

func NewFileStore() (*FileStore, error) {

	var conf FileStoreConfig
	if err := config.ReadConfigEnv(&conf); err != nil {
		return nil, errors.Wrap(err, "failed to read file store config")
	}

	if err := conf.validate(); err != nil {
		return nil, errors.Wrap(err, "failed to validate file store config")
	}

	client, err := connect(conf.Uri)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to mongodb")
	}

	return &FileStore{
		client: client,
		config: conf,
	}, nil
}

// To connect to mongodb
func connect(dsn string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// mongo.Connect return mongo.Client method
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn))
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to mongodb")
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ping mongodb")
	}

	return client, nil
}

func (fs *FileStore) bucket() (*gridfs.Bucket, error) {
	return gridfs.NewBucket(
		fs.client.Database(fs.config.Database),
		options.GridFSBucket().SetName(name),
	)
}
