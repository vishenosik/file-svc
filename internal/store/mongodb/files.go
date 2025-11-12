package mongodb

import (
	"bytes"
	"context"
	"io"

	"github.com/vishenosik/file-svc/internal/entity"
	"github.com/vishenosik/gocherry/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (fs *FileStore) Save(filename string, file []byte) (string, error) {

	bucket, err := fs.bucket()

	if err != nil {
		return "", err
	}

	reader := bytes.NewReader(file)
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, reader); err != nil {
		return "", err
	}

	uploadStream, err := bucket.OpenUploadStream(
		filename,
	)
	if err != nil {
		return "", err
	}
	defer uploadStream.Close()

	if _, err := uploadStream.Write(buf.Bytes()); err != nil {
		return "", err
	}

	id, ok := uploadStream.FileID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("invalid file id")
	}

	return id.Hex(), nil
}

func (fs *FileStore) Get(id string) (file []byte, err error) {

	bucket, err := fs.bucket()
	if err != nil {
		return nil, err
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	if _, err := bucket.DownloadToStream(objID, &buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (fs *FileStore) Delete(id string) error {
	bucket, err := fs.bucket()
	if err != nil {
		return err
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	if err := bucket.Delete(objID); err != nil {
		return err
	}

	return nil
}

type gridfsFile struct {
	Id   string `bson:"_id"`
	Name string `bson:"filename"`
	Size int64  `bson:"chunkSize"`
}

type gridfsFiles = []gridfsFile

func (fs *FileStore) File(id string) (*entity.FileInfo, error) {
	bucket, err := fs.bucket()
	if err != nil {
		return nil, err
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	cursor, err := bucket.Find(bson.M{"_id": objID})
	if err != nil {
		return nil, err
	}

	if !cursor.Next(context.TODO()) {
		return nil, errors.New("file not found")
	}

	var f gridfsFile
	if err := cursor.Decode(&f); err != nil {
		return nil, err
	}

	return fileInfoToEntity(f), nil
}

func (fs *FileStore) Files() (*entity.FileInfoList, error) {
	bucket, err := fs.bucket()
	if err != nil {
		return nil, err
	}

	cursor, err := bucket.Find(bson.D{})
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := cursor.Close(context.TODO()); err != nil {

		}
	}()

	var foundFiles gridfsFiles
	if err = cursor.All(context.TODO(), &foundFiles); err != nil {
		return nil, err
	}

	return filesListToEntity(foundFiles), nil
}

func fileInfoToEntity(f gridfsFile) *entity.FileInfo {
	return &entity.FileInfo{
		ID:   f.Id,
		Name: f.Name,
		Size: f.Size,
	}
}

func filesListToEntity(files gridfsFiles) *entity.FileInfoList {
	list := &entity.FileInfoList{
		Files: make([]*entity.FileInfo, 0, len(files)),
	}

	for _, f := range files {
		list.Files = append(list.Files, fileInfoToEntity(f))
	}
	return list
}
