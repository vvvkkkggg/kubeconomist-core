data "yandex_compute_image" "container-optimized-image" {
  family = "container-optimized-image"
}

resource "yandex_compute_instance" "no-public-ip" {
  count       = 3
  name        = "alkosenko-compute-instance-${count.index + 1}"
  platform_id = "standard-v1"
  zone        = var.zone

  resources {
    cores  = 4
    memory = 8
    core_fraction = 100
  }

  boot_disk {
    initialize_params {
      image_id = data.yandex_compute_image.container-optimized-image.id
    }
  }
  network_interface {
    index     = 1
    subnet_id = "e9brgdgjst6qe6tp1oak"
  }

  metadata = {
    ssh-keys = "ubuntu:${file("~/.ssh/id_ed25519.pub")}"
  }
}

resource "yandex_compute_instance" "with-public-ip" {
    count =  0
  name        = "alkosenko-compute-instance-with-public-ip"
  platform_id = "standard-v2"
  zone        = var.zone

  resources {
    cores  = 8
    memory = 8
    core_fraction = 100
  }

  boot_disk {
    initialize_params {
      image_id = data.yandex_compute_image.container-optimized-image.id
    }
  }

  network_interface {
    subnet_id      = "e9brgdgjst6qe6tp1oak"
    nat            = true
    nat_ip_address = yandex_vpc_address.static_ip[0].external_ipv4_address[0].address
  }

  metadata = {
    ssh-keys = "ubuntu:${file("~/.ssh/id_ed25519.pub")}"
  }
}
