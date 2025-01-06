package process

import (
	"bufio"
	"cloud_storage/mq"
	"cloud_storage/store/ceph"
	"encoding/json"
	"go.uber.org/zap"
	"log"
	"os"
)

// 处理文件转移
func Transfer(msg []byte) bool {
	pubData := mq.TransferData{}
	err := json.Unmarshal(msg, &pubData)
	if err != nil {
		zap.S().Error(err.Error())
		return false
	}

	fin, err := os.Open(pubData.CurLocation)
	if err != nil {
		zap.S().Error(err.Error())
		return false
	}
	cephPath := "/ceph" + pubData.FileHash
	err := ceph.PutObject("userfile", cephPath, bufio.NewReader(fin))
	if err != nil {
		log.Println(err.Error())
		return false
	}

	resp, err := dbcli.UpdateFileLocation(
		pubData.FileHash,
		pubData.DestLocation)
	if err != nil {
		zap.S().Error(err.Error())
		return false
	}

	if !resp.Suc {
		zap.S().Error("更新数据库异常，请检查:" + pubData.FileHash)
		return false
	}
	return true
}
