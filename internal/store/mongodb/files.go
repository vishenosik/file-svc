package mongodb

import (
	"bytes"
	"io"

	"github.com/vishenosik/gocherry/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (fs *FileStore) Save(filename string, file []byte) (string, error) {

	bucket, err := gridfs.NewBucket(
		fs.client.Database(fs.config.Database),
		options.GridFSBucket().SetName(name),
	)

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

	bucket, err := gridfs.NewBucket(
		fs.client.Database(fs.config.Database),
		options.GridFSBucket().SetName(name),
	)
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
