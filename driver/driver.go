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

// Package driver provides the CSI driver
package driver

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vultr/govultr/v3"
	"github.com/vultr/metadata"
	"golang.org/x/oauth2"
	"k8s.io/mount-utils"
	"k8s.io/utils/exec"
)

const (
	DefaultDriverName = "block.csi.vultr.com"
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
	mountID         string

	isController bool
	waitTimeout  time.Duration

	log *logrus.Entry

	mounter *mount.SafeFormatAndMount
	resizer *mount.ResizeFs

	version string
}

func NewDriver(endpoint, token, driverName, version, userAgent, apiURL string) (*VultrDriver, error) {
	if driverName == "" {
		driverName = DefaultDriverName
	}

	ctx := context.Background()
	config := &oauth2.Config{}
	ts := config.TokenSource(ctx, &oauth2.Token{AccessToken: token})
	client := govultr.NewClient(oauth2.NewClient(ctx, ts))

	client.UserAgent = "csi-vultr/" + version

	if userAgent != "" {
		client.UserAgent = fmt.Sprintf("csi-vultr/%s/%s", version, userAgent)
	} else {
		client.UserAgent = "csi-vultr/" + version
	}

	if apiURL != "" {
		if err := client.SetBaseURL(apiURL); err != nil {
			return nil, err
		}
	}

	c := metadata.NewClient()
	meta, err := c.Metadata()
	if err != nil {
		return nil, err
	}

	log := logrus.New().WithFields(logrus.Fields{
		"region":  meta.Region.RegionCode,
		"host_id": meta.InstanceV2ID,
		"version": version,
	})

	return &VultrDriver{
		name:     driverName,
		endpoint: endpoint,
		nodeID:   meta.InstanceV2ID,
		region:   meta.Region.RegionCode,
		client:   client,

		isController: token != "",
		waitTimeout:  defaultTimeout,

		log: log,
		mounter: &mount.SafeFormatAndMount{
			Interface: mount.New(""),
			Exec:      exec.New(),
		},

		resizer: mount.NewResizeFs(mount.SafeFormatAndMount{
			Interface: mount.New(""),
			Exec:      exec.New(),
		}.Exec),

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
