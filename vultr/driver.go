/*
Copyright 2020 Vultr Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package driver

import (
	"time"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/sirupsen/logrus"
	"github.com/vultr/govultr"
)

const (
	defaultTimeout = 1 * time.Minute
)

type VultrDriver struct {
	name          string
	vendorVersion string
	bsPrefix      string
	endpoint      string
	waitTimeout   time.Duration
	log           *logrus.Entry

	identity   csi.IdentityServer
	node       csi.ControllerServer
	controller csi.NodeServer

	account   govultr.AccountService
	snapshot  govultr.SnapshotService
	bsStorage govultr.BlockStorageService
	server    govultr.ServerService
}

func GetDriver() *VultrDriver {
	return &VultrDriver{}
}

// NewDriver returns a CSI plugin that contains gRPC interfaces
// which interact with Kubernetes over unix domain sockets for managing Block Storage
func NewDriver(vultrClient *govultr.Client, driverName, version, prefix, url, ep string) (*VultrDriver, error) {
	driver := GetDriver()

	if driverName == "" {
		driverName = driver.name
	}

	if version == "" {
		version = "dev" // TODO add default version
	}

	// Authenticate client

	// Initialize metadata

	// TODO fix: Set up logging
	log := logrus.New().WithFields(logrus.Fields{
		"region":  "region",
		"host_id": "hostID",
	})

	return &VultrDriver{
		name:          driverName,
		vendorVersion: version,
		endpoint:      ep,
		log:           log,

		// TODO Differentiate driver's purpose: Node or Controller
		// isController:      "",
		waitTimeout: defaultTimeout,

		bsStorage: vultrClient.BlockStorage,
		server:    vultrClient.Server,
		snapshot:  vultrClient.Snapshot,
		account:   vultrClient.Account,
	}, nil
}

func Run() {
	// TODO
}
