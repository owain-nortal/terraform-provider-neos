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

resource "neos_user" "peek-c2" {
  first_name = "peek"
  last_name  = "c1"
  username   = "peekc1"
  email      = "peekc1@peek.org.uk"
  enabled    = true
  account    = "test03"
}

resource "neos_user" "peek-c2a" {
  first_name = "peek"
  last_name  = "c1"
  username   = "peekc1"
  email      = "peekc1@peek.org.uk"
  enabled    = true
  account    = "test02"
}


# resource "neos_user" "peek-b5" {
#   first_name = "peek"
#   last_name  = "five"
#   username   = "peek5"
#   email      = "peek5@peek.org.uk"
#   enabled    = true
#   account    = "test03"
# }

# resource "neos_user" "peek-b1" {
#   first_name = "peek"
#   last_name  = "one"
#   username   = "peeka1"
#   email      = "peeka1@peek.org.uk"
#   enabled    = true
#   account    = "root"
# }

# resource "neos_user" "peek-b6" {
#   first_name = "peek"
#   last_name  = "five"
#   username   = "peek5"
#   email      = "peek5@peek.org.uk"
#   enabled    = true
#   account    = "test02"
# }
