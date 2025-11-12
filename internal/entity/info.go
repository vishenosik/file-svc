package entity

type FileInfo struct {
	ID   string
	Name string
	Size int64
}

type FileInfoList struct {
	Files []*FileInfo
}
