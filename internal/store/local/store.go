package local

type FileStore struct {
}

func NewFileStore() *FileStore {
	return &FileStore{}
}

func (fs *FileStore) Save(name string, file []byte) (id string, err error) {
	return "", nil
}
