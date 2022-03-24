## Installation

### Requirements

[`allow_privileged`](https://www.nomadproject.io/docs/drivers/docker#allow_privileged)
must be enabled in Nomad config.

### Components

To make vultr csi work you need to run two components:

- csi-controller
- csi-node

csi-controller and csi-node uses the same binary, but different
[schedulers](https://www.nomadproject.io/docs/schedulers). csi-controller should
be run as _service_ and csi-node should be run as _system_ (run on every host).

See more at Nomad documentation on CSI
[here](https://www.nomadproject.io/docs/internals/plugins/csi).

You will need to run a separate deployment for each Vultr region.

### API key

In order for the csi to work properly, you will need to provide an API key to
csi-controller.

To obtain a API key, please visit
[API settings](https://my.vultr.com/settings/#settingsapi).

API key can be passed to csi-controller securely via
[Vault integration](https://www.nomadproject.io/docs/integrations/vault-integration).

Example snippet:

```hcl
      task "csi-controller" {
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
```

### Run Vultr CSI

Adapt and run example jobs definitions:

- [controller](examples/csi-controller.nomad.hcl)
- [node](examples/csi-node.nomad.hcl)

In Nomad UI in Storage tab make sure plugin is healthy.

### Create and register example volume

Nomad will not create volume on demand. You need to create a volume yourself
either by hand or with
[Terraform](https://registry.terraform.io/providers/vultr/vultr/latest/docs/resources/block_storage)
and then register it in Nomad, again by hand with
[`nomad volume create`](https://www.nomadproject.io/docs/commands/volume/create)
command or with
[Terraform](https://registry.terraform.io/providers/hashicorp/nomad/latest/docs/resources/volume).

Adapt and use [this](examples/volume.tf) config to test.

### Validate

To validate run [example job](examples/example.job.hcl) and the following
commands:

```shell
nomad exec -job example touch /data/example
nomad stop -purge example
nomad system gc
nomad run example.nomad.hcl
nomad exec -job example ls -alh /data
```

## Examples

Examples of Nomad jobs and Terraform configs can be found [here](examples/).
