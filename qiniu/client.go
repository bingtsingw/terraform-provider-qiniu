package qiniu

import (
	"github.com/bingtsingw/terraform-provider-qiniu/qiniu/sdk/cert"
	"github.com/qiniu/go-sdk/v7/storage"
)

type Client struct {
	bucketconn *storage.BucketManager
	certconn   *cert.CertManager
}
