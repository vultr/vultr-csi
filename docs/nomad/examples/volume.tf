terraform {
  required_providers {
    vultr = {
      source = "vultr/vultr"
    }
  }
}

provider "nomad" {}

data "nomad_plugin" "vultr" {
  plugin_id        = "vultr-ams"
  wait_for_healthy = true
}

resource "vultr_block_storage" "example" {
  lifecycle {
    ignore_changes = [
      attached_to_instance,
    ]
  }

  size_gb = 10
  label   = "example"
  live    = true
  region  = "ams"
}

resource "nomad_volume" "example" {
  type        = "csi"
  namespace   = "default"
  plugin_id   = data.nomad_plugin.vultr.id
  volume_id   = "example"
  name        = "example"
  external_id = vultr_block_storage.example.id

  capability {
    access_mode     = "single-node-writer"
    attachment_mode = "file-system"
  }

  mount_options {
    fs_type     = "ext4"
    mount_flags = ["noatime"]
  }
}
