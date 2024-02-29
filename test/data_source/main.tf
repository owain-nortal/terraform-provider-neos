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

resource "neos_data_source" "op-test1" {
  name        = "datasource-test-1"
  description = "desc test data source 1 updated"
  owner       = "test data source 1 owner"
  label       = "DS1"
  links       = var.links
  contact_ids = var.contact_ids
}
