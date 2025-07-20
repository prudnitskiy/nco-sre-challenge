variable "kubernetes_config_path" {
  type        = string
  default     = "~/.kube/config"
  description = "Path to Kubernetes config file"
}

variable "kubernetes_context" {
  type        = string
  default     = "default"
  description = "Kubernetes context name"
}

variable "app_version" {
  type        = string
  description = "App version to deploy"
  default     = "0.0.1"
}

variable "ingress_host" {
  type        = string
  default     = "nco.prudnitskiy.pro"
  description = "Host service will be available on"
}

variable "app_namespace" {
  type        = string
  default     = "nco"
  description = "K8s namespace app will be deployed to"

}
