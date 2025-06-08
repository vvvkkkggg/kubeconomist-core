# resource "yandex_kubernetes_cluster" "k8s" {
#   name       = "demo-k8s"
#   network_id = var.network_id
#   master {
#     zonal {
#       zone      = var.zone
#       subnet_id = var.subnet_id
#     }
#     public_ip = true
#   }
#   service_account_id      = var.k8s_sa_id
#   node_service_account_id = var.k8s_node_sa_id
#   release_channel         = "RAPID"
#   network_policy_provider = "CALICO"
# }

resource "yandex_kubernetes_node_group" "node_group" {
  cluster_id = "catjf4lqjrdc4u4cbepj"
  name       = "k8s-demo"
  version    = "1.32"
  instance_template {
    platform_id = "standard-v1"
    resources {
      memory = 4
      cores  = 4
    }
    boot_disk {
      size = 64
    }
    network_interface {
      subnet_ids = ["e9brgdgjst6qe6tp1oak"]
      nat        = true
    }
  }
  scale_policy {
    fixed_scale {
      size = 1
    }
  }
  allocation_policy {
    location {
      zone = var.zone
    }
  }
}

# resource "yandex_kubernetes_addon" "prometheus" {
#   cluster_id = yandex_kubernetes_cluster.k8s.id
#   name       = "prometheus"
#   version    = "v2"
#   manifest {
#     content = file("${path.module}/manifests/prometheus.yaml")
#   }
# }

resource "yandex_container_registry" "registry" {
  name = "alkosenko-registry"
  folder_id = var.folder_id
}
