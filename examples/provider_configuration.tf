terraform {
  required_providers {
    cancom = {
      source  = "cancom/cancom"
      version = "0.0.1-pre"
    }
  }
}

provider "cancom" {
  token = "<token>"
}
