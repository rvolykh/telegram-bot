resource "null_resource" "build_binary" {
  triggers = {
    always_run = timestamp()
  }

  provisioner "local-exec" {
    working_dir = "${path.module}/../../../apps/alerting"
    environment = {
      GOOS        = "linux"
      GOARCH      = "amd64"
      CGO_ENABLED = "0"
    }
    command = "go build -trimpath -ldflags=\"-s -w\" -o bootstrap"
  }
}

data "archive_file" "lambda_zip" {
  type        = "zip"
  source_file = "${path.module}/../../../apps/alerting/bootstrap"
  output_path = "${abspath(path.module)}/${var.name}.zip"

  depends_on = [
    // Wait for the binary to be built
    null_resource.build_binary,
  ]
}
