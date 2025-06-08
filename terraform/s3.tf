resource "yandex_storage_bucket" "bucket" {
  count  = 3
  bucket = "alkosenko-test-${count.index + 1}"

  access_key = yandex_iam_service_account_static_access_key.sa-static-key.access_key
  secret_key = yandex_iam_service_account_static_access_key.sa-static-key.secret_key
}

resource "yandex_iam_service_account" "sa" {
  name = "s3"
}

resource "yandex_resourcemanager_folder_iam_member" "sa-admin" {
  folder_id = var.folder_id
  role      = "storage.admin"
  member    = "serviceAccount:${yandex_iam_service_account.sa.id}"
}

resource "yandex_iam_service_account_static_access_key" "sa-static-key" {
  service_account_id = yandex_iam_service_account.sa.id
  description        = "static access key for object storage"
}

resource "null_resource" "generate_file" {
  count = 3
  provisioner "local-exec" {
    command = "head -c 104857600 </dev/urandom > ./dummy_100mb_${count.index + 1}.bin"
  }
}

resource "yandex_storage_object" "test-object" {
  count      = 3
  access_key = yandex_iam_service_account_static_access_key.sa-static-key.access_key
  secret_key = yandex_iam_service_account_static_access_key.sa-static-key.secret_key
  bucket     = "alkosenko-test-${count.index + 1}"
  key        = "test-${count.index + 1}"
  source     = "./dummy_100mb_${count.index + 1}.bin"

  depends_on = [
    yandex_storage_bucket.bucket,
    null_resource.generate_file,
  ]
}
