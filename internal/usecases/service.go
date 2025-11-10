package usecases

func (fs *FileService) GetBatchSize() uint32 {
	return fs.config.BatchSize
}

func (fs *FileService) GetMaxFileSize() uint32 {
	return fs.config.MaxSize
}
