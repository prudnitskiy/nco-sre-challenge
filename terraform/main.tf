locals {
  app_chart_values = {
    image = {
      tag = var.app_version
    }
    ingress = {
      enabled = true
      hosts = [
        {
          host = var.ingress_host
          paths = [
            {
              path     = "/"
              pathType = "Prefix"
            }
          ]
        }
      ]
      tls = [
        {
          hosts = [
            var.ingress_host,
          ]
        }
      ]
    }
  }
}

resource "kubernetes_namespace" "nco_ns" {
  metadata {
    name = var.app_namespace
  }
}

resource "helm_release" "nco_app" {
  name      = "nco-app"
  namespace = kubernetes_namespace.nco_ns.metadata[0].name
  chart     = "${path.root}/../chart/app"
  version   = "0.1.0"
  values    = [yamlencode(local.app_chart_values)]
}
