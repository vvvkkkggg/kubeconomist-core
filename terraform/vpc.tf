resource "yandex_vpc_address" "static_ip" {
  name = "alkosenko-static-ip"
  external_ipv4_address {
    zone_id = var.zone
  }
}

resource "yandex_dns_zone" "zone" {
  name        = "alkosenko-test"
  zone        = "alkosenko.test."
  public      = true
  description = "Demo domain zone"
}
