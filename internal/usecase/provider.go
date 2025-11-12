package usecase

import (
	"log/slog"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/vishenosik/gocherry/pkg/logs"
)

type Provider interface {
	Get(id string) (file []byte, err error)
	Save(name string, file []byte) (id string, err error)
	Delete(id string) error
}

type provider struct {
	log    *slog.Logger
	source Provider
}

func NewProvider(source Provider) *provider {
	return &provider{
		source: source,
		log:    logs.SetupLogger().With(logs.AppComponent("usecase.provider")),
	}
}

func (fs *provider) Upload(filename string, file []byte) (string, error) {

	uid := uuid.New().String()
	ext := filepath.Ext(filename)

	id, err := fs.source.Save(uid+ext, file)
	if err != nil {
		return "", err
	}
	fs.log.Info("file uploaded", slog.String("id", id))
	return id, nil
}

func (fs *provider) Download(id string) (file []byte, err error) {
	file, err = fs.source.Get(id)
	if err != nil {
		return nil, err
	}
	fs.log.Info("file downloaded", slog.String("id", id))
	return file, nil
}

func (fs *provider) DeleteFile(id string) error {
	err := fs.source.Delete(id)
	if err != nil {
		return err
	}
	fs.log.Warn("file deleted", slog.String("id", id))
	return nil
}
