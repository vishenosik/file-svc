package usecase

type Settings interface {
	BatchSize() uint32
	MaxFileSize() uint32
}

type service struct {
	batchSize uint32
	maxSize   uint32
}

func NewService(settings Settings) *service {
	return &service{
		batchSize: settings.BatchSize(),
		maxSize:   settings.MaxFileSize(),
	}
}

func (ss *service) GetBatchSize() uint32 {
	return ss.batchSize
}

func (ss *service) GetMaxFileSize() uint32 {
	return ss.maxSize
}
