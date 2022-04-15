# Change Log

## [v0.6.0](https://github.com/vultr/vultr-csi) (2022-04-15)
### Enhancements 
* Added support for multi-block [71](https://github.com/vultr/vultr-csi/pull/71)

### Dependencies
* GoVultr bump to v2.15.1 [71](https://github.com/vultr/vultr-csi/pull/71)

### Documentation
* Nomad documentation [66](https://github.com/vultr/vultr-csi/pull/66)
* Block Types information [71](https://github.com/vultr/vultr-csi/pull/71)

### Dependencies
* Bump google.golang.org/grpc from 1.40.0 to 1.44.0 [58](https://github.com/vultr/vultr-csi/pull/58)
* Bump google.golang.org/grpc from 1.44.0 to 1.45.0 [61](https://github.com/vultr/vultr-csi/pull/61)
* GoVultr bump to v2.14.1 & fixed FakeInstance [60](https://github.com/vultr/vultr-csi/pull/60)

## [v0.5.0](https://github.com/vultr/vultr-csi) (2022-03-11)
### Dependencies
* Bump google.golang.org/grpc from 1.40.0 to 1.44.0 [58](https://github.com/vultr/vultr-csi/pull/58)
* Bump google.golang.org/grpc from 1.44.0 to 1.45.0 [61](https://github.com/vultr/vultr-csi/pull/61)
* GoVultr bump to v2.14.1 & fixed FakeInstance [60](https://github.com/vultr/vultr-csi/pull/60)

## [v0.4.0](https://github.com/vultr/vultr-csi) (2022-01-19)
### Enhancements
* Update CSIDriver Kind to use API v1 1.22 support [52](https://github.com/vultr/vultr-csi/pull/52)

## [v0.3.0](https://github.com/vultr/vultr-csi) (2021-09-24)
### Dependencies
* Updated all quay images [48](https://github.com/vultr/vultr-csi/pull/48)
* Bumped Go from 1.15 to 1.16 [48](https://github.com/vultr/vultr-csi/pull/48)


## [v0.2.0](https://github.com/vultr/vultr-csi) (2021-06-29)
### Dependencies
* Updated all quay images to pull from GCR + updated their versions [45](https://github.com/vultr/vultr-csi/pull/45)

### Enhancements
* Ability to set custom useragent [43](https://github.com/vultr/vultr-csi/pull/43)

## [v0.1.1](https://github.com/vultr/vultr-csi) (2021-03-25)
### Dependencies
* Update vultr/metadata to v1.0.3 [38](https://github.com/vultr/vultr-csi/pull/38)


## [v0.1.0](https://github.com/vultr/vultr-csi) (2021-03-01)
### Enhancements
* Update to use API v2 [33](https://github.com/vultr/vultr-csi/pull/33)
* Update CSI deps [34](https://github.com/vultr/vultr-csi/pull/34)
* Update to use `mountID` from Vultr API v2 to identify mount path [36](https://github.com/vultr/vultr-csi/pull/36)

## [v0.0.4](https://github.com/vultr/vultr-csi) (2020-11-12)
### Bug Fixes
*  default socket location had wrong path [31](https://github.com/vultr/vultr-csi/pull/31)

### Enhancements
* Cleaned up naming on kubernetes resources to be more uniformed [31](https://github.com/vultr/vultr-csi/pull/31)

### Docker Image
[CSI Container v0.0.4](https://hub.docker.com/r/vultr/vultr-csi/tags)


## [v0.0.3](https://github.com/vultr/vultr-csi) (2020-04-29)
### Dependencies
*  quay.io/k8scsi/csi-attacher v1.0.0 -> v2.2.0[#29](https://github.com/vultr/vultr-csi/pull/29)
*  quay.io/k8scsi/driver-registrar v1.0-canary -> v2.2.0[#29](https://github.com/vultr/vultr-csi/pull/29)
*  quay.io/k8scsi/csi-attacher v1.0.0 -> quay.io/k8scsi/csi-node-driver-registrar:v1.3.0 [#29](https://github.com/vultr/vultr-csi/pull/29)
*  quay.io/k8scsi/csi-provisioner v1.0.0 -> v1.6.0[#29](https://github.com/vultr/vultr-csi/pull/29)
*  goVultr v0.3.2 -> v4.2.0[#29](https://github.com/vultr/vultr-csi/pull/29)
*  metadata v1.0.0 -> v1.0.1[#29](https://github.com/vultr/vultr-csi/pull/29)

### Docker Image
[CSI Container v0.0.3](https://hub.docker.com/layers/vultr/vultr-csi/v0.0.3/images/sha256-1b1b12d4b6b5baab4f3db7f44cbd5055aaa463c84c0bf37d7a2c0aae5201a185?context=explore)

## [v0.0.2](https://github.com/vultr/vultr-csi) (2020-04-29)
### Enhancement
*  Added in vultr metadata client to retrieve information on boot[#26](https://github.com/vultr/vultr-csi/pull/26)

### Docker Image
[CSI Container v0.0.2](https://hub.docker.com/layers/vultr/vultr-csi/v0.0.2/images/sha256-bf31b1d0c92a8af3fc26d67f24ace41cab853f8baeec225e18487259bd7147a8?context=explore)

## [v0.0.1](https://github.com/vultr/vultr-csi) (2020-04-02)

### Initial Release

### Docker Image
[CSI Container v0.0.1](https://hub.docker.com/layers/vultr/vultr-csi/v0.0.1/images/sha256-bddb7d5dbb0ab999f6cb1b34f38036854ed3ca861be2fafdd3d7caadf61b0a53?context=explore)
