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

resource "neos_user" "test-user-a11" {
  first_name = "testa11_first"
  last_name  = "testa11_last"
  username   = "test.usera11"
  email      = "test.usera11@example.com"
  enabled    = true
}


resource "neos_user_policy" "test-user-a11-policy" {
  id = neos_user.test-user-a11.id
  policy_json  = <<EOH
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
            "principal": "${neos_user.test-user-a11.id}",
            "resource": [
              "nrn:ksa:iam::root:account:root"
            ],
            "sid": "root-membership"
          }
        ],
        "version": "2022-10-01"
      },
      "user": "${neos_user.test-user-a11.id}"
    }
  EOH  
}

# data "neos_user" "users" {
# }

output "users" {
  value = neos_user.test-user-a11.id
}

# data "neos_user_policy" "policies" {
# }


# output "data_unit" {
#   value = data.neos_user_policy.policies.policy[0].policy
# }
