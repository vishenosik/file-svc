package config

import (
	"github.com/vishenosik/gocherry/pkg/config"
	"github.com/vishenosik/gocherry/pkg/errors"
	"github.com/vishenosik/gocherry/pkg/logs"
)

func init() {
	config.AddStructs(FileServiceConfig{})
}

type FileServiceConfig struct {
	BatchSize uint32 `env:"FILE_BATCH_SIZE" env-default:"100000"`
	MaxSize   uint32 `env:"FILE_MAX_SIZE" env-default:"10485760"`
}

func (FileServiceConfig) Desc() string {
	return "File service config"
}

func (c FileServiceConfig) Validate() error {
	return nil
}

type settings struct {
	batchSize uint32
	maxSize   uint32
}

func NewService() (*settings, error) {

	log := logs.SetupLogger().With(
		logs.Operation("config.NewService"),
	)

	var conf FileServiceConfig
	if err := config.ReadConfigEnv(&conf); err != nil {
		log.Warn("failed to read env config", logs.Error(err))
	}

	if err := conf.Validate(); err != nil {
		return nil, errors.Wrap(err, "failed to validate authentication service config")
	}

	return &settings{
		batchSize: conf.BatchSize,
		maxSize:   conf.MaxSize,
	}, nil
}

func (s *settings) BatchSize() uint32 {
	return s.batchSize
}

func (s *settings) MaxFileSize() uint32 {
	return s.maxSize
}
