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

resource "neos_data_product" "dp-json" {
  name        = "dp-bike-report-d91111"
  description = "The KPI on bike usage"
  owner       = "Richard C level"
  label       = "DPB"
  links       = [""]
  contact_ids = ["jane C level", "martin C level"]
  builder_json = ""
  schema = {}
}
