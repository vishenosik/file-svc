package usecases

func (fs *FileService) GetBatchSize() uint32 {
	return fs.config.UploadBatchSize
}

func (fs *FileService) GetMaxFileSize() uint32 {
	return fs.config.MaxFileSize
}
