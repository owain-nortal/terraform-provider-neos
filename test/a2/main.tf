terraform {
  required_providers {
    neos = {
      source  = "owain-nortal/neos"
      version = "0.1.1"
    }
  }
}

provider "neos" {
  username  = "neosadmin"
  password  = "**11"
  hub_host  = "owain10.neosdata.cloud"
  core_host = "owain10.neosdata.cloud"
  account   = "root"
  partition = "ksa"
}



resource "neos_account" "account-test03" {
  name         = "test03"
  description  = "description of account "
  owner        = "dave"
  display_name = "account_displayname_1"
}

resource "neos_user" "test-test03" {
  first_name = "test03_first"
  last_name  = "test03_last"
  username   = "test03.user"
  email      = "test03.user@example.com"
  account    = neos_account.account-test03.name
  enabled    = true
}

resource "neos_account" "account-test04" {
  name         = "test04"
  description  = "description of account "
  owner        = "dave"
  display_name = "account_displayname_2"
}

resource "neos_user" "test-test04" {
  first_name = "test03_first"
  last_name  = "test03_last"
  username   = "test03.user"
  email      = "test03.user@example.com"
  account    = neos_account.account-test04.name
  enabled    = true
}
