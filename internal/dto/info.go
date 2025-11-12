package dto

import (
	"github.com/vishenosik/file-svc-sdk/api"
	"github.com/vishenosik/file-svc/internal/entity"
	"github.com/vishenosik/file-svc/internal/usecase"
)

type infoDto struct {
	usc *usecase.Info
}

func NewInfoDTO(usc *usecase.Info) *infoDto {
	return &infoDto{
		usc: usc,
	}
}

func (i *infoDto) GetFileInfo(id string) (*api.FileInfo, error) {
	info, err := i.usc.GetFileInfo(id)
	if err != nil {
		return nil, err
	}

	return fileInfoToApi(info), nil
}

func (i *infoDto) ListFiles() (*api.FileInfoList, error) {
	list, err := i.usc.ListFiles()
	if err != nil {
		return nil, err
	}

	return fileInfoListToApi(list), nil
}

func fileInfoToApi(fi *entity.FileInfo) *api.FileInfo {
	return &api.FileInfo{
		ID:       fi.ID,
		Size:     uint32(fi.Size),
		Filename: fi.Name,
	}
}

func fileInfoListToApi(list *entity.FileInfoList) *api.FileInfoList {

	files := make([]*api.FileInfo, 0, len(list.Files))
	for _, info := range list.Files {
		files = append(files, fileInfoToApi(info))
	}

	return &api.FileInfoList{
		Total: uint32(len(list.Files)),
		Files: files,
	}
}
