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
	_ "context"
	"fmt"
	"log"
	"net"
	_ "net/http"
	"net/url"
	"path"
	"path/filepath"
	"time"

	_ "github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/vultr/govultr"
	"google.golang.org/grpc"
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
	isController  bool

	grpc *grpc.Server

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
	driver.name = driverName
	driver.vendorVersion = version

	if driverName == "" {
		log.Fatalln("Vultr Driver name is missing")
	}

	if version == "" {
		driver.vendorVersion = "dev"
	}

	// TODO metadata

	return &VultrDriver{
		name:          driverName,
		vendorVersion: version,
		endpoint:      ep,
		isController:  false, // TODO Differentiate driver's purpose: Node or Controller
		waitTimeout:   defaultTimeout,

		bsStorage: vultrClient.BlockStorage,
		server:    vultrClient.Server,
		snapshot:  vultrClient.Snapshot,
		account:   vultrClient.Account,
	}, nil
}

// Run starts the pluginwhich will run on the given port
func (driver *VultrDriver) Run(endpoint string) error {
	// test
	fmt.Printf("Running on this endpoint: %s", endpoint)

	// Parse endpoint
	u, err := url.Parse(endpoint)

	if err != nil {
		log.Fatalln(err)
	}

	grpcAddr := path.Join(u.Host, filepath.FromSlash(u.Path))
	if u.Host == "" {
		grpcAddr = filepath.FromSlash(u.Path)
	}

	// Set up gRCP listener
	grpcListener, err := net.Listen(u.Scheme, grpcAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// Register identity server
	// driver.grpc = grpc.NewServer(grpc.UnaryInterceptor(errHandler))
	// csi.RegisterIdentityServer(driver.grpc, driver)

	return nil
}
