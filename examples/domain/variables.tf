variable "domain_name" {}

variable "type" {
  default = "normal"
}

variable "platform" {
  default = "web"
}

variable "geo_cover" {
  default = "global"
}

variable "protocol" {
  default = "http"
}

variable "qiniu_bucket" {}
