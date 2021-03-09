package qiniu

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema:         map[string]*schema.Schema{},
		DataSourcesMap: map[string]*schema.Resource{},
		ResourcesMap:   map[string]*schema.Resource{},
	}

	return provider
}
