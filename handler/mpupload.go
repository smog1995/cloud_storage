package handler

import (
	rPool "cloud_storage/cache/redis"
	mydb "cloud_storage/db"
	"cloud_storage/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

// MultipartUploadInfo : 初始化信息
type MultipartUploadInfo struct {
	FileHash   string
	FileSize   int
	UploadID   string
	ChunkSize  int
	ChunkCount int
}

// InitialMultipartUploadHandler : 初始化分块上传
func InitialMultipartUploadHandler(c *gin.Context) {
	username := c.PostForm("username")
	filehash := c.PostForm("filehash")
	filesize, err := strconv.Atoi(c.PostForm("filesize"))
	if err != nil {
		resp := util.NewRespMsg(-1, "params invalid", nil)
		c.JSON(http.StatusOK, resp.JSONBytes())
		return
	}

	//获取redis连接池
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()
	upInfo := MultipartUploadInfo{
		FileHash:   filehash,
		FileSize:   filesize,
		UploadID:   username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize:  5 * (1 << 20),
		ChunkCount: int(math.Ceil(float64(filesize) / (5 * (1 << 20)))),
	}
	//将初始化信息写入redis缓存
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "chunkcount", upInfo.ChunkCount)
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "filehash", upInfo.FileHash)
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "filesize", upInfo.FileSize)

	c.JSON(http.StatusOK, util.NewRespMsg(0, "OK", upInfo))
}

// UploadPartHandler : 上传文件分块
func UploadPartHandler(c *gin.Context) {
	uploadID := c.PostForm("uploadid")
	chunkIndex := c.PostForm("index")

	//获取redis一个连接
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	fpath := "/data/" + uploadID + "/" + chunkIndex
	os.MkdirAll(path.Dir(fpath), 0744)
	fd, err := os.Create(fpath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.NewRespMsg(-1, "upload part failed", nil).JSONBytes())
		return
	}
	defer fd.Close()

	buf := make([]byte, 1<<20)
	for {
		n, err := c.Request.Body.Read(buf)
		if err != nil {
			break
		}
		fd.Write(buf[:n])

	}

	// 更新redis缓存状态
	rConn.Do("HSET", "MP_"+uploadID, "chkidx_"+chunkIndex, 1)

	c.JSON(http.StatusOK, util.NewRespMsg(0, "OK", nil).JSONBytes())
}

// CompleteUploadHandler : 验证是否上传完毕所有分块，是则通知
func CompleteUploadHandler(c *gin.Context) {

	uploadID := c.PostForm("uploadid")
	username := c.PostForm("username")
	filehash := c.PostForm("filehash")
	filesize := c.PostForm("filesize")
	filename := c.PostForm("filename")

	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	//uploadid查询redis并判断是否所有分块上传完成
	data, err := redis.Values(rConn.Do("HGETALL", "MP_"+uploadID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.NewRespMsg(-1, "complete upload failed", nil).JSONBytes())
		return
	}
	totalCount := 0
	chunkCount := 0
	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i-1].([]byte))
		if k == "chunkcount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			chunkCount++
		}
	}
	if totalCount != chunkCount {
		c.JSON(http.StatusInternalServerError, util.NewRespMsg(-1, "上传失败：分块数量有误", nil).JSONBytes())
		return
	}

	// 合并分块
	fsize, _ := strconv.Atoi(filesize)
	mydb.OnFileUploadFinished(filehash, filename, int64(fsize), "")
	mydb.OnUserFileUploadFinished(username, filehash, filename, int64(fsize))

	c.JSON(http.StatusOK, util.NewRespMsg(0, "OK", nil).JSONBytes())
}
