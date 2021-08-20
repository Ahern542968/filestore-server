package meta

import "filestore-server/db"

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

func UpdateFileMeta(fmeta FileMeta) bool {
	//FileMetas[fmeta.FileSha1] = fmeta
	return db.OnFileUploadFinished(fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)
}

func GetFileMeta(fileSha1 string) FileMeta {
	return FileMetas[fileSha1]
}

func RemoveFileMeta(fileSha1 string) {
	delete(FileMetas, fileSha1)
}