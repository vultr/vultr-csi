job "vultr-csi-nodes" {
  datacenters = ["dc1"]
  namespace   = "default"
  type        = "system"

  group "vultr-ams" {
    task "plugin" {
      driver = "docker"

      config {
        image = "vultr/vultr-csi:v0.5.0"

        privileged = true

        args = [
          "-endpoint=unix:///csi/csi.sock",
        ]
      }

      csi_plugin {
        id        = "vultr-amd"
        type      = "node"
        mount_dir = "/csi"
      }

      resources {
        cpu    = 100
        memory = 64
      }
    }
  }
}
