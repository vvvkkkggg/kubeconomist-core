resource "yandex_storage_bucket" "bucket" {
  bucket     = "alkosenko-test"
  access_key = var.access_key
  secret_key = var.secret_key
  acl        = "private"
}
