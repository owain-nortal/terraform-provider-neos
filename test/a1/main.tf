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

resource "neos_user" "test-user-a15" {
  first_name = "testa15_first"
  last_name  = "testa15_last"
  username   = "test.usera15"
  email      = "test.usera15@example.com"
  enabled    = true
  account    = "root"
}


resource "neos_user_policy" "test-user-a15-policy" {
  id          = neos_user.test-user-a15.id
  policy_json = <<EOH
    {
      "is_system": false,
      "policy": {
        "statements": [
          {
            "action": [
              "account:member",
              "core:announce",             
              "principal:browse"
            ],
            "condition": [],
            "effect": "allow",
            "principal": "${neos_user.test-user-a15.id}",
            "resource": [
              "nrn:ksa:iam::root:account:root"
            ],
            "sid": "root-membership"
          }
        ],
        "version": "2022-10-01"
      },
      "user": "${neos_user.test-user-a15.id}"
    }
  EOH  
}

# data "neos_user" "users" {
# }

output "users" {
  value = neos_user.test-user-a15.id
}

# data "neos_user_policy" "policies" {
# }


# output "data_unit" {
#   value = data.neos_user_policy.policies.policy[0].policy
# }


resource "neos_account" "account-test01" {
  name         = "test01"
  description  = "description of account "
  owner        = "dave"
  display_name = "account_displayname_1"
}

resource "neos_user" "test-test01" {
  first_name = "test01_first"
  last_name  = "test01_last"
  username   = "test01.user"
  email      = "test01.user@example.com"
  account    = neos_account.account-test01.name
  enabled    = true
}

resource "neos_account" "account-test02" {
  name         = "test02"
  description  = "description of account "
  owner        = "dave"
  display_name = "account_displayname_2"
}

resource "neos_user" "test-test02" {
  first_name = "test01_first"
  last_name  = "test01_last"
  username   = "test01.user"
  email      = "test01.user@example.com"
  account    = neos_account.account-test02.name
  enabled    = true
}
