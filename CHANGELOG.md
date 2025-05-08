# Change Log
## [v0.16.0](https://github.com/vultr/vultr-csi/compare/v0.15.1...v0.16.0) (2025-05-08)
### Bug Fixes
* Check mount and device during node staging/attach [PR 317](https://github.com/vultr/vultr-csi/pull/317)

### Dependencies
* Update github-tag-actions from v1.36.0 to v1.71.0 [PR 307](https://github.com/vultr/vultr-csi/pull/307)
* Bump google.golang.org/grpc from 1.71.0 to 1.72.0 [PR 310](https://github.com/vultr/vultr-csi/pull/310)
* Bump golang.org/x/net from 0.36.0 to 0.38.0 [PR 309](https://github.com/vultr/vultr-csi/pull/309)
* Bump k8s.io/mount-utils from 0.32.2 to 0.32.3 [PR 297](https://github.com/vultr/vultr-csi/pull/297)
* Bump golang.org/x/oauth2 from 0.28.0 to 0.29.0 [PR 306](https://github.com/vultr/vultr-csi/pull/306)
* Bump github.com/vultr/govultr/v3 from 3.18.0 to 3.19.1 [PR 312](https://github.com/vultr/vultr-csi/pull/312)
* Bump golang.org/x/oauth2 from 0.29.0 to 0.30.0 [PR 315](https://github.com/vultr/vultr-csi/pull/315)
* Bump github.com/vultr/govultr/v3 from 3.19.1 to 3.20.0 [PR 316](https://github.com/vultr/vultr-csi/pull/316)
* Bump golang.org/x/sys from 0.31.0 to 0.33.0 [PR 314](https://github.com/vultr/vultr-csi/pull/314)
* Bump golang.org/x/sync from 0.12.0 to 0.14.0 [PR 313](https://github.com/vultr/vultr-csi/pull/313)

## [v0.15.1](https://github.com/vultr/vultr-csi/compare/v0.15.0...v0.15.1) (2025-03-28)
### Bug Fixes
* Set default block storage type during NodeStageVolume [PR 300](https://github.com/vultr/vultr-csi/pull/300)

### Dependencies
* Bump golang.org/x/net from 0.34.0 to 0.36.0 [PR 298](https://github.com/vultr/vultr-csi/pull/298)
* Update govultr from v3.16.1 to v3.18.0 [PR 301](https://github.com/vultr/vultr-csi/pull/301)

## [v0.15.0](https://github.com/vultr/vultr-csi/compare/v0.14.0...v0.15.0) (2025-03-12)
### Enhancements
* Add example VFS PVC manifest 

### Automation
* Update workflows to use golangci-lint-action [PR 294](https://github.com/vultr/vultr-csi/pull/294)

### Dependencies
* Update Go version from v1.23 to v1.24 [PR 294](https://github.com/vultr/vultr-csi/pull/294)
* Update govultr from v3.13.0 to v3.16.1 [PR 293](https://github.com/vultr/vultr-csi/pull/293)
* Bump golang.org/x/sync from 0.10.0 to 0.12.0 [PR 291](https://github.com/vultr/vultr-csi/pull/291)
* Bump golang.org/x/oauth2 from 0.25.0 to 0.28.0 [PR 290](https://github.com/vultr/vultr-csi/pull/290)
* Bump google.golang.org/grpc from 1.69.4 to 1.71.0 [PR 289](https://github.com/vultr/vultr-csi/pull/289)
* Bump k8s.io/mount-utils from 0.32.0 to 0.32.2 [PR 284](https://github.com/vultr/vultr-csi/pull/284)

## [v0.14.0](https://github.com/vultr/vultr-csi/compare/v0.13.4...v0.14.0) (2025-01-20)
### Enhancements
* Breaking change: Add support for virtual file system alongside block storage [PR 275](https://github.com/vultr/vultr-csi/pull/275)

### Documentation
* Fix broken link in README.md [PR 267](https://github.com/vultr/vultr-csi/pull/267)

### Automation
* Add builds to public container registry repo [PR 268](https://github.com/vultr/vultr-csi/pull/268)

### Dependencies
* Update govultr from v3.11.2 to v3.13.0 [PR 269](https://github.com/vultr/vultr-csi/pull/269)
* Bump golang.org/x/oauth2 from 0.23.0 to 0.24.0 [PR 257](https://github.com/vultr/vultr-csi/pull/257)
* Bump google.golang.org/grpc from 1.67.1 to 1.69.2 [PR 270](https://github.com/vultr/vultr-csi/pull/270)
* Bump k8s.io/mount-utils from 0.31.2 to 0.32.0 [PR 265](https://github.com/vultr/vultr-csi/pull/265)
* Bump google.golang.org/grpc from 1.69.2 to 1.69.4 [PR 276](https://github.com/vultr/vultr-csi/pull/276)
* Bump golang.org/x/sys from 0.26.0 to 0.29.0 [PR 273](https://github.com/vultr/vultr-csi/pull/273)
* Bump golang.org/x/oauth2 from 0.24.0 to 0.25.0 [PR 272](https://github.com/vultr/vultr-csi/pull/272)
* Bump golang.org/x/sync from 0.8.0 to 0.10.0 [PR 261](https://github.com/vultr/vultr-csi/pull/261)
* Bump github.com/container-storage-interface/spec from 1.10.0 to 1.11.0 [PR 271](https://github.com/vultr/vultr-csi/pull/271)
* Update golang.org/x/net from v0.30.0 to v.33.0 [PR 279](https://github.com/vultr/vultr-csi/pull/279)

### New Contributors
* @JoeyC-Dev made their first contribution in [PR 267](https://github.com/vultr/vultr-csi/pull/267)

## [v0.13.3](https://github.com/vultr/vultr-csi/compare/v0.13.2...v0.13.3) (2024-11-01)
### Bug Fixes
* Fix volume unpublish errors when missing directory [PR 226](https://github.com/vultr/vultr-csi/pull/226)

### Dependencies
* Update govultr from v3.11.1 to v3.11.2 [PR 254](https://github.com/vultr/vultr-csi/pull/254)
* Bump google.golang.org/grpc from 1.66.1 to 1.67.1 [PR 245](https://github.com/vultr/vultr-csi/pull/245)
* Bump golang.org/x/sys from 0.25.0 to 0.26.0 [PR 246](https://github.com/vultr/vultr-csi/pull/246)
* Bump k8s.io/mount-utils from 0.31.0 to 0.31.2 [PR 249](https://github.com/vultr/vultr-csi/pull/249)
* Update govultr from v3.9.1 to v3.11.1 [PR 252](https://github.com/vultr/vultr-csi/pull/252)
* Update go from v1.22 to v1.23 [PR 251](https://github.com/vultr/vultr-csi/pull/251)
* Upgrade to Go 1.22 and Update Vendor Dependencies [PR 242](https://github.com/vultr/vultr-csi/pull/242)

### New Contributors
* @beornf made their first contribution in [PR 226](https://github.com/vultr/vultr-csi/pull/226)

## [v0.13.2](https://github.com/vultr/vultr-csi/compare/v0.13.1...v0.13.2) (2024-08-12)
### Automation
* Fix releaser mattermost notification action [PR 221](https://github.com/vultr/vultr-csi/pull/221)

### Dependencies
* Bump golang.org/x/sync from 0.6.0 to 0.8.0 [PR 231](https://github.com/vultr/vultr-csi/pull/231)
* Bump golang.org/x/oauth2 from 0.20.0 to 0.22.0 [PR 230](https://github.com/vultr/vultr-csi/pull/230)
* Bump google.golang.org/grpc from 1.64.0 to 1.65.0 [PR 228](https://github.com/vultr/vultr-csi/pull/228)

## [v0.13.1](https://github.com/vultr/vultr-csi/compare/v0.13.0...v0.13.1) (2024-06-17)
### Automation
* Update goreleaser workflow to use action v6 and go v1.22 [PR 219](https://github.com/vultr/vultr-csi/pull/219)

## [v0.13.0](https://github.com/vultr/vultr-csi/compare/v0.12.4...v0.13.0) (2024-06-17)
### Dependencies
* Bump golang.org/x/oauth2 from 0.17.0 to 0.20.0 [PR 210](https://github.com/vultr/vultr-csi/pull/210)
* Bump github.com/vultr/govultr/v3 from 3.6.2 to 3.8.1 [PR 215](https://github.com/vultr/vultr-csi/pull/215)
* Update golang/protobuf from 1.5.3 to 1.5.4 [PR 216](https://github.com/vultr/vultr-csi/pull/216)
* Bump golang.org/x/net from 0.21.0 to 0.23.0 [PR 208](https://github.com/vultr/vultr-csi/pull/208)
* Bump golang.org/x/sys from 0.20.0 to 0.21.0 [PR 217](https://github.com/vultr/vultr-csi/pull/217)
* Bump google.golang.org/grpc from 1.61.1 to 1.64.0 [PR 214](https://github.com/vultr/vultr-csi/pull/214)

## [v0.12.4](https://github.com/vultr/vultr-csi/compare/v0.12.3...v0.12.4) (2024-05-13)
### Bug Fix
* Set container image to resolve XFS mounting issues [PR 212](https://github.com/vultr/vultr-csi/pull/212)

### Documentation
* Update region availability in docs [PR 203](https://github.com/vultr/vultr-csi/pull/203)

### Automation
* Update mattermost notify action [PR 209](https://github.com/vultr/vultr-csi/pull/209)

## [v0.12.3](https://github.com/vultr/vultr-csi/compare/v0.12.2...v0.12.3) (2024-02-21)
### Bug Fix
* Fix ARM builds with type conversion [PR 193](https://github.com/vultr/vultr-csi/pull/193)

## [v0.12.2](https://github.com/vultr/vultr-csi/compare/v0.12.1...v0.12.2) (2024-02-21)
### Bug Fix
* Resolve CSI expansion issues [PR 189](https://github.com/vultr/vultr-csi/pull/189)
* Fix some lint issues and update the linter configuration [PR 191](https://github.com/vultr/vultr-csi/pull/191)

### Dependencies
* Bump k8s.io/mount-utils from 0.29.1 to 0.29.2 [PR 188](https://github.com/vultr/vultr-csi/pull/188)
* Bump google.golang.org/grpc from 1.60.1 to 1.61.1 [PR 187](https://github.com/vultr/vultr-csi/pull/187)
* Bump golang.org/x/oauth2 from 0.16.0 to 0.17.0 [PR 186](https://github.com/vultr/vultr-csi/pull/186)
* Bump github.com/vultr/govultr/v3 from 3.6.1 to 3.6.2 [PR 190](https://github.com/vultr/vultr-csi/pull/190)

## [v0.12.1](https://github.com/vultr/vultr-csi/compare/v0.12.0...v0.12.1) (2024-02-04)
### Enhancements
* Add node volume resize [PR 184](https://github.com/vultr/vultr-csi/pull/184)

## [v0.12.0](https://github.com/vultr/vultr-csi/compare/v0.11.0...v0.12.0) (2024-01-23)
### Bug Fix
* Fix mount invalid argument [PR 180](https://github.com/vultr/vultr-csi/pull/180)

### Dependencies
* Bump golang.org/x/sys from 0.14.0 to 0.16.0 [PR 173](https://github.com/vultr/vultr-csi/pull/173)
* Bump google.golang.org/grpc from 1.59.0 to 1.60.1 [PR 172](https://github.com/vultr/vultr-csi/pull/172)
* Bump golang.org/x/oauth2 from 0.14.0 to 0.16.0 [PR 175](https://github.com/vultr/vultr-csi/pull/175)
* Bump golang.org/x/sync from 0.4.0 to 0.6.0 [PR 174](https://github.com/vultr/vultr-csi/pull/174)

### Automation
* Use GITHUB_OUTPUT envvar instead of set-output command [PR 176](https://github.com/vultr/vultr-csi/pull/176)

### New Contributors
* @arunsathiya made their first contribution in [PR 176](https://github.com/vultr/vultr-csi/pull/176)

## [v0.11.0](https://github.com/vultr/vultr-csi/compare/v0.10.1...v0.11.0) (2023-11-16)
### Bug Fixes
* Decrease node maximum mount limit [PR 162](https://github.com/vultr/vultr-csi/pull/162)

### Documentation
* Tidy up some changelog links [PR 150](https://github.com/vultr/vultr-csi/pull/150)

### Dependencies
* Bump google.golang.org/grpc from 1.58.2 to 1.59.0 [PR 154](https://github.com/vultr/vultr-csi/pull/154)
* Bump github.com/vultr/govultr/v3 from 3.3.1 to 3.4.0 [PR 160](https://github.com/vultr/vultr-csi/pull/160)
* Bump golang.org/x/oauth2 from 0.12.0 to 0.14.0 [PR 159](https://github.com/vultr/vultr-csi/pull/159)
* Bump github.com/container-storage-interface/spec from 1.8.0 to 1.9.0 [PR 155](https://github.com/vultr/vultr-csi/pull/155)

## [v0.10.1](https://github.com/vultr/vultr-csi/compare/v0.10.0...v0.10.1) (2023-10-06)
### Automation
* Resolve build errors from statfs type inconsistencies in darwin arm64 [PR 148](https://github.com/vultr/vultr-csi/pull/148)

### Dependencies
* Bump golang.org/x/sys from 0.12.0 to 0.13.0 [PR 147](https://github.com/vultr/vultr-csi/pull/147)
* Bump golang.org/x/sync from 0.3.0 to 0.4.0 [PR 146](https://github.com/vultr/vultr-csi/pull/146)

## [v0.10.0](https://github.com/vultr/vultr-csi/compare/v0.9.0...v0.10.0) (2023-10-05)
### Enhancements
* Update NVMe mininum block size [PR 143](https://github.com/vultr/vultr-csi/pull/143)

### Bug Fixes
* Fix formatting error handling [PR 119](https://github.com/vultr/vultr-csi/pull/119)

### Documentation
* Fix typo in csi_plugin.id in csi-node.nomad.hcl example [PR 124](https://github.com/vultr/vultr-csi/pull/124)

### Automation
* Replace golint with golangci-lint in go-checks workflow [PR 123](https://github.com/vultr/vultr-csi/pull/123)

### Dependencies
* Bump google.golang.org/grpc from 1.52.3 to 1.56.1 [PR 125](https://github.com/vultr/vultr-csi/pull/125)
* Update Go to v1.20 [PR 126](https://github.com/vultr/vultr-csi/pull/126)
* Bump google.golang.org/grpc from 1.56.1 to 1.56.2 [PR 127](https://github.com/vultr/vultr-csi/pull/127)
* Bump golang.org/x/oauth2 from 0.7.0 to 0.10.0 [PR 129](https://github.com/vultr/vultr-csi/pull/129)
* Bump golang.org/x/sync from 0.0.0-20210220032951-036812b2e83c to 0.3.0 [PR 128](https://github.com/vultr/vultr-csi/pull/128)
* Update govultr to v3.1.0 [PR 132](https://github.com/vultr/vultr-csi/pull/132)
* Bump golang.org/x/oauth2 from 0.10.0 to 0.12.0 [PR 138](https://github.com/vultr/vultr-csi/pull/138)
* Bump github.com/sirupsen/logrus from 1.9.0 to 1.9.3 [PR 131](https://github.com/vultr/vultr-csi/pull/131)
* Update govultr to v3.3.1 [PR 144](https://github.com/vultr/vultr-csi/pull/144)
* Bump github.com/container-storage-interface/spec from 1.7.0 to 1.8.0 [PR 130](https://github.com/vultr/vultr-csi/pull/130)
* Bump google.golang.org/grpc from 1.56.2 to 1.58.2 [PR 141](https://github.com/vultr/vultr-csi/pull/141)

### New Contributors
* @const-tmp made their first contribution in [PR 124](https://github.com/vultr/vultr-csi/pull/124)
* @mondragonfx made their first contribution in [PR 143](https://github.com/vultr/vultr-csi/pull/143)

## [v0.9.0](https://github.com/vultr/vultr-csi/compare/v0.8.0...v0.9.0) (2023-03-06)
### Enhancements
* Added volume expansion capability in [PR 116](https://github.com/vultr/vultr-csi/pull/116)
* Added volume statistics capability in [PR 115](https://github.com/vultr/vultr-csi/pull/115)

## [v0.8.0](https://github.com/vultr/vultr-csi/compare/v0.7.0...v0.8.0) (2023-01-30)
### Dependencies
* Bump github.com/sirupsen/logrus v1.9.0
* Bump github.com/vultr/govultr/v2 v2.17.2
* Bump google.golang.org/grpc v1.52.3
* Bump github.com/container-storage-interface/spec v1.7.0

## [v0.7.0](https://github.com/vultr/vultr-csi/compare/v0.6.0...v0.7.0) (2022-05-13)
### Enhancements
* Allow changing Vultr API url [PR 75](https://github.com/vultr/vultr-csi/pull/75)

### Dependencies
* Bump github.com/container-storage-interface/spec from 1.5.0 to 1.6.0 [PR 69](https://github.com/vultr/vultr-csi/pull/69)
* Bump github.com/vultr/metadata from 1.0.3 to 1.1.0 [PR 76](https://github.com/vultr/vultr-csi/pull/76)
* Bump google.golang.org/grpc from 1.45.0 to 1.46.0 [PR 73](https://github.com/vultr/vultr-csi/pull/73)
* Bump GO to 1.17 [PR 77](https://github.com/vultr/vultr-csi/pull/77)
* Bump github.com/vultr/govultr/v2 from 2.15.1 to 2.16.0 [PR 74](https://github.com/vultr/vultr-csi/pull/74)


## [v0.6.0](https://github.com/vultr/vultr-csi/compare/v0.5.0...v0.6.0) (2022-04-15)
### Enhancements
* Added support for multi-block [71](https://github.com/vultr/vultr-csi/pull/71)

### Dependencies
* GoVultr bump to v2.15.1 [71](https://github.com/vultr/vultr-csi/pull/71)

### Documentation
* Nomad documentation [66](https://github.com/vultr/vultr-csi/pull/66)
* Block Types information [71](https://github.com/vultr/vultr-csi/pull/71)

## [v0.5.0](https://github.com/vultr/vultr-csi/compare/v0.4.0...v0.5.0) (2022-03-11)
### Dependencies
* Bump google.golang.org/grpc from 1.40.0 to 1.44.0 [58](https://github.com/vultr/vultr-csi/pull/58)
* Bump google.golang.org/grpc from 1.44.0 to 1.45.0 [61](https://github.com/vultr/vultr-csi/pull/61)
* GoVultr bump to v2.14.1 & fixed FakeInstance [60](https://github.com/vultr/vultr-csi/pull/60)

## [v0.4.0](https://github.com/vultr/vultr-csi/compare/v0.3.0...v0.4.0) (2022-01-19)
### Enhancements
* Update CSIDriver Kind to use API v1 1.22 support [52](https://github.com/vultr/vultr-csi/pull/52)

## [v0.3.0](https://github.com/vultr/vultr-csi/compare/v0.2.0...v0.3.0) (2021-09-24)
### Dependencies
* Updated all quay images [48](https://github.com/vultr/vultr-csi/pull/48)
* Bumped Go from 1.15 to 1.16 [48](https://github.com/vultr/vultr-csi/pull/48)


## [v0.2.0](https://github.com/vultr/vultr-csi/compare/v0.1.1...v0.2.0) (2021-06-29)
### Dependencies
* Updated all quay images to pull from GCR + updated their versions [45](https://github.com/vultr/vultr-csi/pull/45)

### Enhancements
* Ability to set custom useragent [43](https://github.com/vultr/vultr-csi/pull/43)

## [v0.1.1](https://github.com/vultr/vultr-csi/compare/v0.1.0...v0.1.1) (2021-03-25)
### Dependencies
* Update vultr/metadata to v1.0.3 [38](https://github.com/vultr/vultr-csi/pull/38)


## [v0.1.0](https://github.com/vultr/vultr-csi/compare/v0.0.4...v0.1.0) (2021-03-01)
### Enhancements
* Update to use API v2 [33](https://github.com/vultr/vultr-csi/pull/33)
* Update CSI deps [34](https://github.com/vultr/vultr-csi/pull/34)
* Update to use `mountID` from Vultr API v2 to identify mount path [36](https://github.com/vultr/vultr-csi/pull/36)

## [v0.0.4](https://github.com/vultr/vultr-csi/compare/v0.0.3...v0.0.4) (2020-11-12)
### Bug Fixes
*  default socket location had wrong path [31](https://github.com/vultr/vultr-csi/pull/31)

### Enhancements
* Cleaned up naming on kubernetes resources to be more uniformed [31](https://github.com/vultr/vultr-csi/pull/31)

### Docker Image
[CSI Container v0.0.4](https://hub.docker.com/r/vultr/vultr-csi/tags)


## [v0.0.3](https://github.com/vultr/vultr-csi/compare/v0.0.2...v0.0.3) (2020-04-29)
### Dependencies
*  quay.io/k8scsi/csi-attacher v1.0.0 -> v2.2.0[#29](https://github.com/vultr/vultr-csi/pull/29)
*  quay.io/k8scsi/driver-registrar v1.0-canary -> v2.2.0[#29](https://github.com/vultr/vultr-csi/pull/29)
*  quay.io/k8scsi/csi-attacher v1.0.0 -> quay.io/k8scsi/csi-node-driver-registrar:v1.3.0 [#29](https://github.com/vultr/vultr-csi/pull/29)
*  quay.io/k8scsi/csi-provisioner v1.0.0 -> v1.6.0[#29](https://github.com/vultr/vultr-csi/pull/29)
*  goVultr v0.3.2 -> v4.2.0[#29](https://github.com/vultr/vultr-csi/pull/29)
*  metadata v1.0.0 -> v1.0.1[#29](https://github.com/vultr/vultr-csi/pull/29)

### Docker Image
[CSI Container v0.0.3](https://hub.docker.com/layers/vultr/vultr-csi/v0.0.3/images/sha256-1b1b12d4b6b5baab4f3db7f44cbd5055aaa463c84c0bf37d7a2c0aae5201a185?context=explore)

## [v0.0.2](https://github.com/vultr/vultr-csi/release/tag/v0.0.2) (2020-04-29)
### Enhancement
*  Added in vultr metadata client to retrieve information on boot[#26](https://github.com/vultr/vultr-csi/pull/26)

### Docker Image
[CSI Container v0.0.2](https://hub.docker.com/layers/vultr/vultr-csi/v0.0.2/images/sha256-bf31b1d0c92a8af3fc26d67f24ace41cab853f8baeec225e18487259bd7147a8?context=explore)

## [v0.0.1](https://github.com/vultr/vultr-csi) (2020-04-02)

### Initial Release

### Docker Image
[CSI Container v0.0.1](https://hub.docker.com/layers/vultr/vultr-csi/v0.0.1/images/sha256-bddb7d5dbb0ab999f6cb1b34f38036854ed3ca861be2fafdd3d7caadf61b0a53?context=explore)
