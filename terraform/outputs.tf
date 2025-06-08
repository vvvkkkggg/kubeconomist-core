output "static_ip" {
  value = yandex_vpc_address.static_ip.external_ipv4_address[0].address
}

output "bucket_name" {
  value = yandex_storage_bucket.bucket.bucket
}

output "registry_id" {
  value = yandex_container_registry.registry.id
}

output "k8s_cluster_id" {
  value = yandex_kubernetes_cluster.k8s.id
}

output "dns_zone_id" {
  value = yandex_dns_zone.zone.id
}
