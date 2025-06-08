resource "yandex_vpc_address" "static_ip" {
  count = 3
  name  = "alkosenko-static-ip-${count.index + 1}"
  external_ipv4_address {
    zone_id = var.zone
  }
}

resource "yandex_dns_zone" "zone" {
  count       = 5
  name        = "alkosenko-test-${count.index + 1}"
  zone        = "alkosenko${count.index + 1}.test."
  public      = true
  description = "Demo domain zone ${count.index + 1}"
}
