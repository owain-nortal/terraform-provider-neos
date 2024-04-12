terraform {
  required_providers {
    neos = {
      source = "registry.terraform.io/owain-nortal/neos"
    }
  }
}

provider "neos" {
  username  = "neosadmin"
  password  = "**"
  hub_host  = "owain10.neosdata.cloud"
  core_host = "owain10.neosdata.cloud"
  account   = "root"
  partition = "ksa"
}

variable "links" {
  type    = list(any)
  default = ["link1", "link2"]
}

variable "contact_ids" {
  type    = list(any)
  default = ["contacts1", "contacts2"]
}

resource "neos_data_source" "op-test3" {
  name            = "datasource-test-3"
  description     = "desc test data source 1 updated"
  owner           = "test data source 1 owner"
  label           = "DS1"
  links           = var.links
  contact_ids     = var.contact_ids
  connection_json = jsonencode({ "connection" : { "access_key" : { "env_key" : "MINIO_S3_ACCESS" }, "access_secret" : { "env_key" : "MINIO_S3_SECRET" }, "connection_type" : "s3", "url" : "http://minio.neos-core-minio.svc.cluster.local" } })
  secret_values = {
    MINIO_S3_ACCESS = "abc123"
    MINIO_S3_SECRET = "bar"
  }
}


resource "neos_data_source" "op-test4" {
  name            = "datasource-test-4"
  description     = "desc test data source 1 updated"
  owner           = "test data source 1 owner"
  label           = "DS1"
  links           = var.links
  contact_ids     = var.contact_ids
  connection_json = jsonencode({ "connection" : { "access_key" : { "env_key" : "MINIO_S3_ACCESS" }, "access_secret" : { "env_key" : "MINIO_S3_SECRET" }, "connection_type" : "s3", "url" : "http://minio.neos-core-minio.svc.cluster.local" } })
  secret_values = {
    MINIO_S3_ACCESS = "secret"
    MINIO_S3_SECRET = "password"
  }
}