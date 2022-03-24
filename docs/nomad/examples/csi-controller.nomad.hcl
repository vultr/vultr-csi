job "vultr-csi-controller" {
  datacenters = ["dc1"]
  namespace   = "default"

  group "vultr-ams" {
    task "controller" {
      driver = "docker"

      vault {
        policies = ["vultr"]
      }

      config {
        image = "vultr/vultr-csi:v0.5.0"

        args = [
          "-endpoint=unix:///csi/csi.sock",
          "-token=${VULTR_API_KEY}",
        ]
      }

      template {
        data = <<-EOF
          VULTR_API_KEY={{ with secret "secret/vultr/csi" }}{{ .Data.data.key }}{{ end }}
          EOF

        destination = "secrets/api.env"
        change_mode = "restart"
        env         = true
      }

      csi_plugin {
        id        = "vultr-ams"
        type      = "controller"
        mount_dir = "/csi"
      }

      resources {
        cpu    = 100
        memory = 64
      }
    }
  }
}
