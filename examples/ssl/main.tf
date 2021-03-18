resource "qiniu_ssl_cert" "test" {
  name = "test"
  ca = file("test.crt")
  pri = file("test.key")
}
