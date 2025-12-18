resource "null_resource" "build_binary" {
  triggers = {
    always_run = timestamp()
  }

  provisioner "local-exec" {
    working_dir = abspath(var.source_path)
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
  source_file = "${abspath(var.source_path)}/bootstrap"
  output_path = "${abspath(path.module)}/${var.function_name}.zip"

  depends_on = [
    // Wait for the binary to be built
    null_resource.build_binary,
  ]
}
