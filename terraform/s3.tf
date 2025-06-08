resource "yandex_storage_bucket" "bucket" {
  count  = 3
  bucket = "alkosenko-test-${count.index + 1}"
  acl    = "private"
}

resource "null_resource" "generate_file" {
  count = 3
  provisioner "local-exec" {
    command = "head -c 104857600 </dev/urandom > ./dummy_100mb_${count.index + 1}.bin"
  }
}

resource "yandex_storage_object" "cute-cat-picture" {
  count = 3

  bucket = "alkosenko-test-${count.index + 1}"
  key    = "test-${count.index + 1}"
  source = "./dummy_100mb_${count.index + 1}.bin"
}
