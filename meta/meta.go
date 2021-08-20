package meta

import (
	"filestore-server/db"
)

type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var FileMetas map[string]FileMeta

func init() {
	FileMetas = make(map[string]FileMeta)
}

func UpdateFileMeta(fmeta FileMeta) {
	FileMetas[fmeta.FileSha1] = fmeta
}

func UpdateFileMetaDB(fmeta FileMeta) bool {
	return db.OnFileUploadFinished(fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)
}

func GetFileMeta(fileSha1 string) FileMeta {
	return FileMetas[fileSha1]
}

func GetFileMetaDB(fileSha1 string) (*FileMeta, error) {
	ftable, err := db.GetFileMeta(fileSha1)
	if err != nil {
		return nil, err
	}
	fmeta := FileMeta{
		FileSha1: ftable.FileSha1.String,
		FileName: ftable.FileName.String,
		FileSize: ftable.FileSize.Int64,
		Location: ftable.FileAddr.String,
	}
	return &fmeta, err
}


func RemoveFileMeta(fileSha1 string) {
	delete(FileMetas, fileSha1)
}