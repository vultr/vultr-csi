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

	"github.com/sirupsen/logrus"
	"github.com/vultr/govultr"
)

const (
	DefaultDriverName = "vultrbs.csi.driver.com"
	defaultTimeout    = 1 * time.Minute
)

// VultrDriver struct
type VultrDriver struct {
	name     string
	endpoint string
	nodeID   string
	region   string
	client   *govultr.Client

	publishVolumeID string

	isController bool
	waitTimeout  time.Duration

	log     *logrus.Entry
	mounter Mounter

	version string
}

func NewDriver(endpoint, token, driverName, version, nodeID string) (*VultrDriver, error) {
	if driverName == "" {
		driverName = DefaultDriverName
	}

	client := govultr.NewClient(nil, token)
	client.UserAgent = "csi-vultr/" + version

	instance, err := GetVultrByName(client, nodeID)
	if err != nil {
		return nil, err
	}

	log := logrus.New().WithFields(logrus.Fields{
		"region":  instance.RegionID,
		"host_id": instance.InstanceID,
		"version": version,
	})

	return &VultrDriver{
		name:     driverName,
		endpoint: endpoint,
		nodeID:   instance.InstanceID,
		region:   instance.RegionID,
		client:   client,

		isController: token != "",
		waitTimeout:  defaultTimeout,

		log:     log,
		mounter: NewMounter(log),

		version: version,
	}, nil
}

func (d *VultrDriver) Run() {
	server := NewNonBlockingGRPCServer()
	identity := NewVultrIdentityServer(d)
	controller := NewVultrControllerServer(d)
	node := NewVultrNodeDriver(d)

	server.Start(d.endpoint, identity, controller, node)
	server.Wait()
}
