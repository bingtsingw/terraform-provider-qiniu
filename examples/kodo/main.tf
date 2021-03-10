data "qiniu_kodo_buckets" "buckets" {
  region_id = "z1"
}

output "buckets" {
  value = data.qiniu_kodo_buckets.buckets
}
