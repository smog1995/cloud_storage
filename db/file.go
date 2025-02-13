package db

import (
	mydb "cloud_storage/db/mysql"
	"database/sql"
	"go.uber.org/zap"
)

// OnFileUploadFinished : 文件上传完成,保存meta
// TODO 这里只实现了插入新行的逻辑，没有实现更新的逻辑
func OnFileUploadFinished(filehash string, filename string, filesize int64, fileaddr string) bool {
	stmt, err := mydb.DBConn().Prepare("insert into ignore into tbl_file" +
		"(`file_sha1`,`file_name`,`file_size`,`file_addr`,`status`) values(?,?,?,?,1)")
	if err != nil {
		zap.S().Error("Failed to prepare statement,err:%s", err.Error())
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		zap.S().Fatalf(err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			zap.S().Error("File with hash been upload before%s", filehash)
		}
		return true
	}
	return false
}

// TableFile : 文件表结构体
type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

// GetFileMeta : 从mysql获取文件元信息
func GetFileMeta(filehash string) (*TableFile, error) {
	stmt, err := mydb.DBConn().Prepare(
		"select file_sha1,file_addr,file_name,file_size from tbl_file" +
			"where file_sha1=? and status=1 limit 1")
	if err != nil {
		zap.S().Fatalf(err.Error())
		return nil, err
	}
	defer stmt.Close()

	tfile := TableFile{}
	err = stmt.QueryRow(filehash).Scan(
		&tfile.FileHash, &tfile.FileAddr, &tfile.FileName, &tfile.FileSize)
	if err != nil {
		if err == sql.ErrNoRows {
			// 查不到对应记录， 返回参数及错误均为nil
			return nil, nil
		}
		zap.S().Fatalf(err.Error())
		return nil, err
	}
	return &tfile, nil
}
