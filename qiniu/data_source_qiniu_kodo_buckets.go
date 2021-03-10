package qiniu

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/qiniu/go-sdk/v7/storage"
	"strconv"
	"time"
)

func dataSourceQiniuKodoBuckets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceQiniuKodoBucketsRead,
		Schema: map[string]*schema.Schema{
			"region_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateRegionID,
				ForceNew:     true,
			},
			"buckets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"index_page_on": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"max_age": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceQiniuKodoBucketsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	conn := meta.(Client).bucketconn
	regionId := storage.RegionID(d.Get("region_id").(string))

	bucketInfos, err := conn.BucketInfosInRegion(regionId, false)
	if err != nil {
		return diag.FromErr(err)
	}

	buckets := make([]map[string]interface{}, 0, len(bucketInfos))
	for _, bucket := range bucketInfos {
		attributes := map[string]interface{}{
			"name":          bucket.Name,
			"region_id":     bucket.Info.Region,
			"private":       bucket.Info.IsPrivate(),
			"index_page_on": bucket.Info.IndexPageOn(),
			"max_age":       bucket.Info.MaxAge,
		}

		buckets = append(buckets, attributes)
	}

	if err := d.Set("buckets", buckets); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
