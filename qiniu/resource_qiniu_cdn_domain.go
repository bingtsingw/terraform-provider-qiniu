package qiniu

import (
	"context"
	"fmt"
	"time"

	"github.com/bingtsingw/terraform-provider-qiniu/qiniu/sdk/domain"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceQiniuCdnDomain() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceQiniuCdnDomainRead,
		CreateContext: resourceQiniuCdnDomainCreate,
		UpdateContext: resourceQiniuCdnDomainUpdate,
		DeleteContext: resourceQiniuCdnDomainDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"normal", "wildcard"}, false),
			},
			"platform": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"web", "download", "vod", "dynamic"}, false),
			},
			"geo_cover": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"china", "foreign", "global"}, false),
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validation.StringInSlice([]string{"http", "https"}, false),
			},
			"https": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cert_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"force": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"http2": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},
			"source": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"domain", "ip", "qiniuBucket", "advanced"}, false),
						},
						// 后台没有该设置选项, [API](https://developer.qiniu.com/fusion/4249/product-features#2)也不清晰
						//"host": {
						//	Type:     schema.TypeString,
						//	Required: true,
						//},
						"ips": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"domain": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"qiniu_bucket": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"url_scheme": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"http", "https"}, false),
						},
						"test_url_path": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "qiniu_do_not_delete.gif",
						},
					},
				},
			},
			"cache": {
				Type:     schema.TypeSet,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ignore_param": {
							Type:     schema.TypeBool,
							Required: true,
							ForceNew: false,
						},
						"controls": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"time": {
										Type:     schema.TypeInt,
										Required: true,
										ForceNew: false,
									},
									"timeunit": {
										Type:         schema.TypeInt,
										Required:     true,
										ForceNew:     false,
										ValidateFunc: validation.IntInSlice([]int{0, 1, 2, 3, 4, 5, 6}),
									},
									"type": {
										Type:         schema.TypeString,
										Required:     true,
										ForceNew:     false,
										ValidateFunc: validation.StringInSlice([]string{"all", "path", "suffix", "follow"}, false),
									},
									"rule": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: false,
									},
								},
							},
						},
					},
				},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
		},
	}
}

func resourceQiniuCdnDomainRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := m.(Client).domainconn
	domainName := d.Id()

	res, err := conn.GetDomainInfo(domainName)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", res.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("type", res.Type); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("platform", res.Platform); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("geo_cover", res.GeoCover); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("protocol", res.Protocol); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("source", flattenResponseDomainSource(res.Source)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("cache", flattenResponseDomainCache(res.Cache)); err != nil {
		return diag.FromErr(err)
	}

	if res.Protocol == "https" {
		if err := d.Set("https", flattenResponseDomainHttps(res.Https)); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceQiniuCdnDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := m.(Client).domainconn

	domainName := d.Get("name").(string)

	input := domain.DomainInfo{
		Type:     d.Get("type").(string),
		Platform: d.Get("platform").(string),
		GeoCover: d.Get("geo_cover").(string),
		Protocol: d.Get("protocol").(string),
		Source:   convertInputDomainSource(d.Get("source").(*schema.Set).List()),
	}

	if protocol, ok := d.GetOk("protocol"); ok {
		if protocol == "https" {
			https := d.Get("https").(*schema.Set).List()
			if len(https) != 1 {
				return diag.FromErr(fmt.Errorf("when protocol is 'https', https block must be set"))
			}

			input.Https = convertInputDomainHttps(https)
		}
	}

	if cache, ok := d.GetOk("cache"); ok {
		input.Cache = convertInputDomainCache(cache.(*schema.Set).List())
	}

	_, err := conn.CreateDomain(domainName, input)

	if err != nil {
		return diag.FromErr(err)
	}

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		res, err := conn.DescribeDomain(domainName)

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error describing domain: %s", err))
		}

		if res.OperationType == "create_domain" && res.OperatingState == "processing" {
			return resource.RetryableError(fmt.Errorf("domain creation is processing"))
		}

		if res.OperationType == "create_domain" && res.OperatingState == "success" {
			return nil
		}

		return resource.NonRetryableError(fmt.Errorf("error describing domain: unkown state"))
	})

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(domainName)

	diags = resourceQiniuCdnDomainRead(ctx, d, m)

	return diags
}

func resourceQiniuCdnDomainUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conn := m.(Client).domainconn
	domainName := d.Id()

	protocol := d.Get("protocol").(string)
	if d.HasChange("protocol") {
		if protocol == "http" {
			// HTTPS降级为HTTP
			err := conn.UnsslizeDomain(domainName)
			if err != nil {
				return diag.FromErr(err)
			}

			err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				res, err := conn.DescribeDomain(domainName)

				if err != nil {
					return resource.NonRetryableError(fmt.Errorf("[unsslize] error describing domain: %s", err))
				}

				if res.OperationType == "unsslize" && res.OperatingState == "processing" {
					return resource.RetryableError(fmt.Errorf("domain unsslize is processing"))
				}

				if res.OperationType == "unsslize" && res.OperatingState == "success" {
					return nil
				}

				return resource.NonRetryableError(fmt.Errorf("[unsslize] error describing domain: unkown state"))
			})

			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			// HTTP升级为HTTPS
			https := convertInputDomainHttps(d.Get("https").(*schema.Set).List())
			err := conn.SslizeDomain(domainName, https)
			if err != nil {
				return diag.FromErr(err)
			}

			err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				res, err := conn.DescribeDomain(domainName)

				if err != nil {
					return resource.NonRetryableError(fmt.Errorf("[sslize] error describing domain: %s", err))
				}

				if res.OperationType == "sslize" && res.OperatingState == "processing" {
					return resource.RetryableError(fmt.Errorf("domain sslize is processing"))
				}

				if res.OperationType == "sslize" && res.OperatingState == "success" {
					return nil
				}

				return resource.NonRetryableError(fmt.Errorf("[sslize] error describing domain: unkown state"))
			})

			if err != nil {
				return diag.FromErr(err)
			}
		}
	} else {
		// 当protocol未改变 && protocol == "https" && https改变
		// 执行证书更新的逻辑
		if protocol == "https" && d.HasChange("https") {
			https := convertInputDomainHttps(d.Get("https").(*schema.Set).List())
			err := conn.ModifyDomainHttpsConf(domainName, https)
			if err != nil {
				return diag.FromErr(err)
			}

			err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				res, err := conn.DescribeDomain(domainName)

				if err != nil {
					return resource.NonRetryableError(fmt.Errorf("[https] error describing domain: %s", err))
				}

				if res.OperationType == "modify_https_conf" && res.OperatingState == "processing" {
					return resource.RetryableError(fmt.Errorf("domain https is processing"))
				}

				if res.OperationType == "modify_https_conf" && res.OperatingState == "success" {
					return nil
				}

				return resource.NonRetryableError(fmt.Errorf("[https] error describing domain: unkown state"))
			})

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceQiniuCdnDomainRead(ctx, d, m)
}

func resourceQiniuCdnDomainDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := m.(Client).domainconn
	domainName := d.Id()

	err := conn.OfflineDomain(domainName)
	if err != nil {
		return diag.FromErr(err)
	}

	err = conn.DeleteDomain(domainName)

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		res, err := conn.DescribeDomain(domainName)

		if err != nil {
			if err.Error() == "无此域名" {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		if res.OperationType == "delete_domain" && res.OperatingState == "processing" {
			return resource.RetryableError(fmt.Errorf("domain creation is processing"))
		}

		return resource.NonRetryableError(fmt.Errorf("error describing domain: unkown state"))
	})

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func flattenResponseDomainHttps(h domain.DomainHttpsInfo) []interface{} {
	https := map[string]interface{}{
		"cert_id": h.CertID,
		"force":   h.ForceHttps,
		"http2":   h.Http2Enable,
	}

	return []interface{}{https}
}

func convertInputDomainHttps(hh []interface{}) domain.DomainHttpsInfo {
	h := hh[0].(map[string]interface{})

	https := domain.DomainHttpsInfo{
		CertID:      h["cert_id"].(string),
		ForceHttps:  h["force"].(bool),
		Http2Enable: h["http2"].(bool),
	}

	return https
}

func flattenResponseDomainSource(s domain.DomainSourceInfo) []interface{} {
	source := map[string]interface{}{
		"type":          s.Type,
		"test_url_path": s.TestURLPath,
	}

	if s.Type == "qiniuBucket" {
		source["qiniu_bucket"] = s.QiniuBucket
	} else {
		source["url_scheme"] = s.URLScheme
	}

	if s.Type == "ip" {
		source["ips"] = s.IPs
	}

	if s.Type == "domain" {
		source["domain"] = s.Domain
	}

	return []interface{}{source}
}

func convertInputDomainSource(ss []interface{}) domain.DomainSourceInfo {
	s := ss[0].(map[string]interface{})
	source := domain.DomainSourceInfo{
		Type:        s["type"].(string),
		IPs:         expandStringList(s["ips"].([]interface{})),
		Domain:      s["domain"].(string),
		QiniuBucket: s["qiniu_bucket"].(string),
		URLScheme:   s["url_scheme"].(string),
		TestURLPath: s["test_url_path"].(string),
	}

	return source
}

func flattenResponseDomainCache(c domain.DomainCacheInfo) []interface{} {
	controls := make([]map[string]interface{}, len(c.CacheControls))

	for i, v := range c.CacheControls {
		control := map[string]interface{}{
			"time":     v.Time,
			"timeunit": v.Timeunit,
			"type":     v.Type,
			"rule":     v.Rule,
		}
		controls[i] = control
	}

	cache := map[string]interface{}{
		"ignore_param": c.IgnoreParam,
		"controls":     controls,
	}

	return []interface{}{cache}
}

func convertInputDomainCache(cc []interface{}) domain.DomainCacheInfo {
	c := cc[0].(map[string]interface{})
	cache := domain.DomainCacheInfo{
		IgnoreParam:   c["ignore_param"].(bool),
		CacheControls: convertInputDomainCacheControls(c["controls"].(*schema.Set).List()),
	}

	return cache
}

func convertInputDomainCacheControls(cc []interface{}) []domain.DomainCacheControl {
	var controls []domain.DomainCacheControl

	for _, c := range cc {
		v := c.(map[string]interface{})
		control := domain.DomainCacheControl{
			Time:     v["time"].(int),
			Timeunit: v["timeunit"].(int),
			Type:     v["type"].(string),
			Rule:     v["rule"].(string),
		}

		controls = append(controls, control)
	}

	return controls
}
