package meta

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

func GetFileMeta(fileSha1 string) FileMeta {
	return FileMetas[fileSha1]
}
