package qiniu

import "github.com/qiniu/go-sdk/v7/storage"

type Client struct {
	bucketconn *storage.BucketManager
}