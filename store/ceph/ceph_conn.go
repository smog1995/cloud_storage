package ceph

import (
	"cloud_storage/global"
	"gopkg.in/amz.v1/aws"
	"gopkg.in/amz.v1/s3"

)

var cephConn *s3.S3

// 获取ceph连接
func GetcephConnection *s3.S3 {
	if cephConn != nil {
		return cephConn
	}

	// 1.初始化ceph信息
	auth := aws.Auth{
		AccessKey: global.ServerConfig.CephAccessKey,
		SecretKey: global.ServerConfig.CephSecretKey,
	}

	curRegion := aws.Region{
		Name: "Default",
		EC2Endpoint: global.ServerConfig.CephGWEndpoint,
		S3Endpoint: global.ServerConfig.CephGWEndpoint,
		S3BucketEndpoint: "",
		S3LocationConstraint: false,
		S3LowercaseBucket: false,
		Sign: aws.SignV2,
	}

	return s3.New(auth, curRegion)
}

// GetCephBucket : 获取指定的bucket对象
func GetCephBucket(bucket string) *s3.Bucket {
	conn := GetCephConnection()
	return conn.Bucket(bucket)
}

// PutObject : 上传文件到ceph集群
func PutObject(bucket string, path string, data []byte) error {
	return GetCephBucket(bucket).Put(path, data, "octet-stream", s3.PublicRead)
}
