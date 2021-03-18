package qiniu

import (
	"github.com/bingtsingw/terraform-provider-qiniu/qiniu/sdk/cert"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
)

type Config struct {
	AccessKey string
	SecretKey string
}

func (c *Config) Client() Client {
	credentials := auth.New(c.AccessKey, c.SecretKey)

	client := Client{
		bucketconn: storage.NewBucketManager(credentials, nil),
		certconn:   cert.NewCertManager(credentials),
	}

	return client
}
