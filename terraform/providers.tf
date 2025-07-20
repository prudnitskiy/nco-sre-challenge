terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 4.37"
    }
    azuread = {
      source  = "hashicorp/azuread"
      version = "~> 3.4"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.17"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 3.0"
    }
  }
  backend "azurerm" {
    resource_group_name  = "tf"
    storage_account_name = "eaitfstate"
    container_name       = "tas"
    key                  = "nco.tfstate"
    use_oidc             = true
  }
}

provider "kubernetes" {
  config_path    = var.kubernetes_config_path
  config_context = var.kubernetes_context
}

provider "helm" {
  kubernetes = {
    config_path    = var.kubernetes_config_path
    config_context = var.kubernetes_context
  }
}
