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

resource "neos_group" "group-test1" {
  name        = "group-test-a8"
  description = "description of group"
  principals = [
    "87fb6e9f-9e97-4a38-b45d-d22d70aef1f4",
    "bee9524a-cd73-4d7f-920d-e65e46a8209a",
    "64760505-298b-4865-b67d-1269c974bf39",
  ]
}
