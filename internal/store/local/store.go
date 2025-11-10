package local

type FileStore struct {
}

func NewFileStore() *FileStore {
	return &FileStore{}
}

func (fs *FileStore) Save(name string, file []byte) (id string, err error) {
	return "", nil
}

func (fs *FileStore) Get(id string) (file []byte, err error) {
	return nil, nil
}
