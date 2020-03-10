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
	// _ "context"
	"fmt"
	"log"
	"net"

	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/sirupsen/logrus"
	"github.com/vultr/govultr"
	"google.golang.org/grpc"
)

const (
	defaultTimeout = 1 * time.Minute
)

// VultrDriver struct
type VultrDriver struct {
	name          string
	vendorVersion string
	bsPrefix      string
	endpoint      string
	waitTimeout   time.Duration
	isController  bool

	log     *logrus.Entry
	grpc    *grpc.Server
	httpSrv *http.Server

	idServer *VultrIdentityServer

	account   govultr.AccountService
	snapshot  govultr.SnapshotService
	bsStorage govultr.BlockStorageService
	server    govultr.ServerService
}

// GetDriver returns VultrDriver
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

	// TODO change to whatever
	log := logrus.New().WithFields(logrus.Fields{
		"version": version,
	})

	return &VultrDriver{
		name:          driverName,
		vendorVersion: version,
		endpoint:      ep,
		isController:  false, // TODO Differentiate driver's purpose: Node or Controller
		waitTimeout:   defaultTimeout,

		log: log,

		bsStorage: vultrClient.BlockStorage,
		server:    vultrClient.Server,
		snapshot:  vultrClient.Snapshot,
		account:   vultrClient.Account,
	}, nil
}

// Run starts the plugin which will run on the given port
func (driver *VultrDriver) Run(endpoint string) error {
	// Parse endpoint and get address
	u, err := url.Parse(endpoint)
	if err != nil {
		log.Fatalln(err)
	}

	grpcAddr := path.Join(u.Host, filepath.FromSlash(u.Path))
	if u.Host == "" {
		grpcAddr = filepath.FromSlash(u.Path)
	}

	// Remove socket if already exists
	if err := os.RemoveAll(grpcAddr); err != nil {
		return fmt.Errorf("could not remove unix domain socket %s, error: %s", grpcAddr, err)
	}

	if u.Scheme != "unix" {
		return fmt.Errorf("only unix domain sockets are supported, have: %s", u.Scheme)
	}

	// Set up gRCP listener
	grpcListener, err := net.Listen(u.Scheme, grpcAddr)
	if err != nil {
		return fmt.Errorf("cannot listen on socket: %v", err)
	}

	// Register identity server
	driver.grpc = grpc.NewServer(grpc.UnaryInterceptor(driver.GRPCLogger))
	// TODO: register other servers here
	csi.RegisterIdentityServer(driver.grpc, driver)

	// test start
	driver.log.WithField("grpc_addr", grpcAddr).Info("server starting...")

	if driver.httpSrv == nil {
		driver.log.WithField("grpc_addr", grpcAddr).Info("server running...")
		return driver.grpc.Serve(grpcListener)
	}

	err = driver.httpSrv.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}
