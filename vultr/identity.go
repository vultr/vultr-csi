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
	"context"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/protobuf/ptypes/wrappers"
)

// VultrIdentityServer ...
type VultrIdentityServer struct {
	Driver *VultrDriver
}

// GetPluginInfo returns basic plugin data
func (d *VultrDriver) GetPluginInfo(ctx context.Context, req *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {
	res := &csi.GetPluginInfoResponse{
		Name:          d.name,
		VendorVersion: d.vendorVersion,
	}
	return res, nil
}

// GetPluginCapabilities returns plugins available capabilities
func (d *VultrDriver) GetPluginCapabilities(ctx context.Context, req *csi.GetPluginCapabilitiesRequest) (*csi.GetPluginCapabilitiesResponse, error) { // This can change depending on what we will offer right away
	volExpType := csi.PluginCapability_VolumeExpansion_ONLINE

	return &csi.GetPluginCapabilitiesResponse{
		Capabilities: []*csi.PluginCapability{
			{
				Type: &csi.PluginCapability_Service_{
					Service: &csi.PluginCapability_Service{
						Type: csi.PluginCapability_Service_CONTROLLER_SERVICE,
					},
				},
			},
			{
				Type: &csi.PluginCapability_VolumeExpansion_{
					VolumeExpansion: &csi.PluginCapability_VolumeExpansion{
						Type: volExpType,
					},
				},
			},
		},
	}, nil
}

// Probe returns plugin health and readiness
func (d *VultrDriver) Probe(ctx context.Context, req *csi.ProbeRequest) (*csi.ProbeResponse, error) {
	return &csi.ProbeResponse{
		Ready: &wrappers.BoolValue{Value: true},
	}, nil
}
