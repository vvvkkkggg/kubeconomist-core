resource "yandex_kubernetes_cluster" "k8s" {
  name        = "demo-k8s"
  network_id  = var.network_id
  master {
    zonal {
      zone      = var.zone
      subnet_id = var.subnet_id
    }
    public_ip = true
  }
  service_account_id      = var.k8s_sa_id
  node_service_account_id = var.k8s_node_sa_id
  release_channel         = "RAPID"
  network_policy_provider = "CALICO"
}

resource "yandex_kubernetes_node_group" "node_group" {
  cluster_id = yandex_kubernetes_cluster.k8s.id
  name       = "k8s-demo"
  version    = "1.28"
  instance_template {
    platform_id = "standard-v3"
    resources {
      memory = 4
      cores  = 2
    }
    boot_disk {
      size = 50
    }
    network_interface {
      subnet_ids = [var.subnet_id]
      nat        = true
    }
  }
  scale_policy {
    fixed_scale {
      size = 2
    }
  }
  allocation_policy {
    location {
      zone = var.zone
    }
  }
  service_account_id = var.k8s_node_sa_id
}

resource "yandex_kubernetes_addon" "prometheus" {
  cluster_id = yandex_kubernetes_cluster.k8s.id
  name       = "prometheus"
  version    = "v2"
  manifest {
    content = file("${path.module}/manifests/prometheus.yaml")
  }
}

resource "yandex_container_registry" "registry" {
  name = "alkosenko-registry"
}
