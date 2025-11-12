package usecase

import (
	"log/slog"

	"github.com/vishenosik/file-svc/internal/entity"
	"github.com/vishenosik/gocherry/pkg/logs"
)

type InfoProvider interface {
	File(id string) (info *entity.FileInfo, err error)
	Files() (list *entity.FileInfoList, err error)
}

type Info struct {
	log    *slog.Logger
	source InfoProvider
}

func NewInfo(source InfoProvider) *Info {
	return &Info{
		source: source,
		log:    logs.SetupLogger().With(logs.AppComponent("usecase.info")),
	}
}

func (i *Info) GetFileInfo(id string) (info *entity.FileInfo, err error) {
	info, err = i.source.File(id)
	if err != nil {
		return nil, err
	}
	i.log.Info("file info retrieved", slog.String("id", id))
	return info, nil
}

func (i *Info) ListFiles() (list *entity.FileInfoList, err error) {
	list, err = i.source.Files()
	if err != nil {
		return nil, err
	}
	i.log.Info("files list retrieved", slog.Int("total", len(list.Files)))
	return list, nil
}
