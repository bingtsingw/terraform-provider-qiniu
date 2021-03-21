package qiniu

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("QINIU_ACCESS_KEY", ""),
			},
			"secret_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("QINIU_SECRET_KEY", ""),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"qiniu_kodo_buckets": dataSourceQiniuKodoBuckets(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"qiniu_ssl_cert":   resourceQiniuSslCert(),
			"qiniu_cdn_domain": resourceQiniuCdnDomain(),
		},
		ConfigureContextFunc: providerConfigure,
	}

	return provider
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := Config{
		AccessKey: d.Get("access_key").(string),
		SecretKey: d.Get("secret_key").(string),
	}

	return config.Client(), nil
}
