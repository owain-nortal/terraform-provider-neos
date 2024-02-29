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



resource "neos_account" "op-test1" {
  name         = "account-test-1"
  description  = "description of account "
  owner        = "dave"
  display_name = "account_displayname_1"
}
