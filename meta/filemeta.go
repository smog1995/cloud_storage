package meta

type FileMeta struct {
	FileSha1 string
	Filename string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// UpdateFilemeta: 新增/更新文件原信息
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

// UpdateFileMetaDB:新增/更新文件元信息到mysql中
func UpdateFileMetaDB(fmeta FileMeta) bool {

}

func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}
