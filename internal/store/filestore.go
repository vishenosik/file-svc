package store

import (
	"log/slog"

	"github.com/pkg/errors"
	"github.com/vishenosik/file-svc/internal/store/local"
	"github.com/vishenosik/file-svc/internal/store/mongodb"
	"github.com/vishenosik/gocherry/pkg/config"
	"github.com/vishenosik/gocherry/pkg/logs"
)

const (
	DriverLocal   = "local"
	DriverMongoDb = "mongodb"
)

func init() {
	config.AddStructs(FileStoreConfig{})
}

type FileStoreConfig struct {
	Driver string `env:"FILE_STORE_DRIVER" default:"local" desc:"file store driver (local, mongodb)"`
}

func (FileStoreConfig) Desc() string {
	return "File storage config"
}

func (c FileStoreConfig) validate() error {
	switch c.Driver {
	case DriverLocal, DriverMongoDb:
	default:
		return errors.Errorf("unknown driver %s", c.Driver)
	}
	return nil
}

type FileSaver interface {
	Save(name string, file []byte) (id string, err error)
}

type FileStore struct {
	log   *slog.Logger
	saver FileSaver

	config FileStoreConfig
}

func NewFileStore() (*FileStore, error) {

	log := logs.SetupLogger().With(
		logs.Operation("NewFileStore"),
	)

	var conf FileStoreConfig
	if err := config.ReadConfigEnv(&conf); err != nil {
		log.Warn("failed to read env config", logs.Error(err))
	}

	if err := conf.validate(); err != nil {
		return nil, errors.Wrap(err, "failed to validate authentication service config")
	}

	var (
		saver FileSaver
		err   error
	)

	switch conf.Driver {
	case DriverLocal:
		saver = local.NewFileStore()
	case DriverMongoDb:
		saver, err = mongodb.NewFileStore()
		if err != nil {
			return nil, errors.Wrap(err, "failed to create file store")
		}
	}

	if saver == nil {
		return nil, errors.New("saver is nil")
	}

	return &FileStore{
		saver:  saver,
		config: conf,
		log:    log,
	}, nil
}

func (fs *FileStore) Save(name string, file []byte) (id string, err error) {
	return fs.saver.Save(name, file)
}
