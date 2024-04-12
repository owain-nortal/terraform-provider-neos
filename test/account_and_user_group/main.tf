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

resource "neos_account" "account-accy1" {
  name         = "accy1"
  description  = "description of account "
  owner        = "dave"
  display_name = "account_displayname_1"
}

resource "neos_user" "test-acc" {
  first_name = "accy_first"
  last_name  = "accy_last"
  username   = "accy.user"
  email      = "accy.user@example.com"
  account    = neos_account.account-accy1.name
  enabled    = true
}

resource "neos_user_policy" "test-acc-policy" {
  id          = neos_user.test-acc.id
  account     = neos_account.account-accy1.name
  policy_json = <<EOH
    {
      "is_system": false,
      "policy": {
        "statements": [
          {
            "action": [
              "account:member",           
              "principal:browse"
            ],
            "condition": [],
            "effect": "allow",
            "principal": "${neos_user.test-acc.id}",
            "resource": [
              "nrn:ksa:iam::root:account:root"
            ],
            "sid": "root-membership"
          }
        ],
        "version": "2022-10-01"
      },
      "user": "${neos_user.test-acc.id}"
    }
  EOH  
}

resource "neos_group" "group-test1" {
  name        = "group-test-a8b"
  description = "description of group"
  principals = [
    "87fb6e9f-9e97-4a38-b45d-d22d70aef1f4",
    "64760505-298b-4865-b67d-1269c974bf39",
  ]
  account = neos_account.account-accy1.name
}
