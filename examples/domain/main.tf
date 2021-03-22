resource "qiniu_cdn_domain" "test" {
  name = var.domain_name
  type = var.type
  platform = var.platform
  geo_cover = var.geo_cover
  protocol = var.protocol

  source {
    type = "qiniuBucket"
    qiniu_bucket = var.qiniu_bucket
  }

  cache {
    ignore_param = false

    controls {
      time = 1
      timeunit = 6
      type = "all"
      rule = "*"
    }
  }
}

output "test" {
  value = qiniu_cdn_domain.test.cname
}
