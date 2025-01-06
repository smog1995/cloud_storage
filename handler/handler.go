package handler

import (
	mydb "cloud_storage/db"
	"cloud_storage/global"
	"cloud_storage/meta"
	"cloud_storage/util"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

// 处理文件上传
func UploadHandler(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "文件上传",
		})

	} else if c.Request.Method == http.MethodPost {
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			zap.S().Error("上传文件失败:  %s", err.Error())
			return
		}
		defer file.Close()

		//文件表的数据行基本信息
		fileMeta := meta.FileMeta{
			FileName: header.Filename,
			Location: global.ServerConfig.FileLocation,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		newFile, err := os.Create(fileMeta.Location)

		if err != nil {
			zap.S().Error("Failed to create file,err:%s\n", err.Error())
			return
		}
		defer newFile.Close()
		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			zap.S().Error("Failed to save data into file,err:%s\n", err.Error())
			return
		}

		//1. 如果文件写入ceph，那么该做的操作如下(location也要更新):
		//newFile.Seek(0, 0)
		//data, _ := ioutil.ReadAll(newFile)
		// bucket:=ceph.GetCephBucket("userfile")
		// cephPath:="/ceph/"+fileMeta.FileSha1
		// _=bucket.Put(cephPath,data,"octet-stream",s3.PublicRead)
		// fileMeta.Location=cephPath

		//2. 如果文件写入oss则是
		//ossPath := "oss/" + fileMeta.FileSha1
		//err = oss.Bucket().PutObject(ossPath, newFile)
		//if err != nil {
		//	fmt.Println(err.Error())
		//	w.Write([]byte("Upload failed!"))
		//	return
		//}
		//fileMeta.Location = ossPath

		fileMeta.FileSha1 = util.FileSha1(newFile)

		// 更新文件的meta
		_ = meta.UpdateFileMetaDB(fileMeta)

		// 更新用户文件表
		username := c.PostForm("username")
	}

}

func GetFileMetaHandler(c *gin.Context) {
	filehash := c.PostForm("filehash")

	fMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		zap.S().Error(err.Error())
		c.String(http.StatusInternalServerError, "server error")
		return
	}
	data, err := json.Marshal(fMeta)
	if err != nil {
		zap.S().Error(err.Error())
		c.String(http.StatusInternalServerError, "server error")
		return
	}
	c.JSON(http.StatusOK, data)

}

// FileQueryHandler : 查询批量的文件元信息
func FileQueryHandler(c *gin.Context) {
	limitCount, _ := strconv.Atoi(c.PostForm("limit"))
	username := c.PostForm("username")
	userFiles, err := mydb.QueryUserFileMetas(username, limitCount)
	if err != nil {
		zap.S().Error("FileQueryHandler: %s", err.Error())
		c.String(http.StatusInternalServerError, "server error")
		return
	}

	data, err = json.Marshal(userFiles)
	if err != nil {
		zap.S().Error("json编码出错: %s", err.Error())
		c.String(http.StatusInternalServerError, "server error")
		return
	}
	c.JSON(http.StatusOK, data)

}

// FileMetaUpdateHandler 更新元信息接口(重命名)
func FileMetaUpdateHandler(c *gin.Context) {
	filesha1 := c.PostForm("filehash")
	newFileName := c.PostForm("filename")
	opType := c.PostForm("op")
	if opType != "0" {
		c.String(http.StatusForbidden, "禁止的行为")
	}
	curFileMeta := meta.GetFileMeta(filesha1)
	curFileMeta.FileName = newFileName
	meta.UpdateFileMeta(curFileMeta)

}
