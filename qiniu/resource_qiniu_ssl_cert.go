package qiniu

import (
	"context"
	"github.com/bingtsingw/terraform-provider-qiniu/qiniu/sdk/cert"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceQiniuSslCert() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceQiniuSslCertCreate,
		ReadContext:   resourceQiniuSslCertRead,
		DeleteContext: resourceQiniuSslCertDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"pri": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ca": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceQiniuSslCertCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := m.(Client).certconn

	c, err := conn.CreateCert(cert.CertBody{
		Name: d.Get("name").(string),
		Pri:  d.Get("pri").(string),
		Ca:   d.Get("ca").(string),
	})

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(c.Id)

	diags = resourceQiniuSslCertRead(ctx, d, m)

	return diags
}

func resourceQiniuSslCertRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := m.(Client).certconn

	c, err := conn.GetCertInfo(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", c.Name); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceQiniuSslCertDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := m.(Client).certconn

	err := conn.DeleteCert(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
