package mq

import (
	cmn "cloud_storage/common"
)

type TransferData struct {
	FileHash      string
	CurLocation   string
	DestLocation  string
	DestStoreType cmn.StoreType
}
