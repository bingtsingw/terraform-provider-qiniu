package qiniu

import (
	"context"
	"github.com/bingtsingw/terraform-provider-qiniu/qiniu/sdk/cert"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceQiniuSslCert() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceQiniuSslCertRead,
		CreateContext: resourceQiniuSslCertCreate,
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
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"ca": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"common_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dns_names": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceQiniuSslCertRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := m.(Client).certconn

	c, err := conn.GetCertInfo(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", c.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("common_name", c.CommonName); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("dns_names", c.DnsNames); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("pri", c.Pri); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ca", c.Ca); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceQiniuSslCertCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := m.(Client).certconn

	c, err := conn.CreateCert(cert.CertInfo{
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

func resourceQiniuSslCertDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := m.(Client).certconn

	err := conn.DeleteCert(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
