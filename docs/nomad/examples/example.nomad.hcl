job "example" {
  datacenters = ["dc1"]
  namespace   = "default"

  group "example" {
    volume "example" {
      type            = "csi"
      access_mode     = "single-node-writer"
      attachment_mode = "file-system"
      source          = "example"
    }

    task "example" {
      driver = "docker"

      volume_mount {
        volume      = "example"
        destination = "/data"
      }

      config {
        image   = "busybox"
        command = "sleep"
        args    = ["1000000"]
      }

      resources {
        cpu    = 10
        memory = 32
      }
    }
  }
}
