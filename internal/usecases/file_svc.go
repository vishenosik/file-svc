package usecases

import (
	"log/slog"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/vishenosik/gocherry/pkg/config"
	"github.com/vishenosik/gocherry/pkg/logs"
)

func init() {
	config.AddStructs(FileServiceConfig{})
}

type FileServiceConfig struct {
}

func (FileServiceConfig) Desc() string {
	return "File service config"
}

func (c FileServiceConfig) validate() error {
	return nil
}

type FileSaver interface {
	Save(name string, file []byte) (id string, err error)
}

type FileService struct {
	log   *slog.Logger
	saver FileSaver

	config FileServiceConfig
}

func NewFileService(
	saver FileSaver,
) (*FileService, error) {

	log := logs.SetupLogger().With(
		logs.Operation("NewFileService"),
	)

	var conf FileServiceConfig
	if err := config.ReadConfigEnv(&conf); err != nil {
		log.Warn("failed to read env config", logs.Error(err))
	}

	if err := conf.validate(); err != nil {
		return nil, errors.Wrap(err, "failed to validate authentication service config")
	}

	if saver == nil {
		return nil, errors.New("saver is nil")
	}

	return &FileService{
		saver:  saver,
		config: conf,
		log:    log,
	}, nil
}

func (fs *FileService) Upload(filename string, file []byte) (string, error) {

	uid := uuid.New().String()
	ext := filepath.Ext(filename)

	id, err := fs.saver.Save(uid+ext, file)
	if err != nil {
		return "", err
	}
	return id, nil
}
