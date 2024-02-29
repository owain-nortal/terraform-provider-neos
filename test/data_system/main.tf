terraform {
  required_providers {
    neos = {
      source = "registry.terraform.io/owain-nortal/neos"
    }
  }
}

provider "neos" {
  username  = "neosadmin"
  password  = "**Gwen11"
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

resource "neos_data_system" "op-test1" {
  name        = "ds-test-1"
  description = "desc test data system 1 updated"
  owner       = "test data system 1 owner"
  label       = "DS1"
  links       = var.links
  contact_ids = var.contact_ids
}
