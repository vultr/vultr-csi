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

var _ csi.IdentityServer = &VultrIdentityServer{}

// VultrIdentityServer
type VultrIdentityServer struct {
	Driver *VultrDriver
}

func NewVultrIdentityServer(driver *VultrDriver) *VultrIdentityServer {
	return &VultrIdentityServer{driver}
}

// GetPluginInfo returns basic plugin data
func (vultrIdentity *VultrIdentityServer) GetPluginInfo(context.Context, *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {
	vultrIdentity.Driver.log.Info("VultrIdentityServer.GetPluginInfo called")

	res := &csi.GetPluginInfoResponse{
		Name:          vultrIdentity.Driver.name,
		VendorVersion: vultrIdentity.Driver.version,
	}
	return res, nil
}

// GetPluginCapabilities returns plugins available capabilities
func (vultrIdentity *VultrIdentityServer) GetPluginCapabilities(_ context.Context, req *csi.GetPluginCapabilitiesRequest) (*csi.GetPluginCapabilitiesResponse, error) {
	vultrIdentity.Driver.log.Infof("VultrIdentityServer.GetPluginCapabilities called with request : %v", req)

	return &csi.GetPluginCapabilitiesResponse{
		Capabilities: []*csi.PluginCapability{
			{
				Type: &csi.PluginCapability_Service_{
					Service: &csi.PluginCapability_Service{
						Type: csi.PluginCapability_Service_CONTROLLER_SERVICE,
					},
				},
			},
		},
	}, nil
}

func (vultrIdentity *VultrIdentityServer) Probe(_ context.Context, req *csi.ProbeRequest) (*csi.ProbeResponse, error) {
	vultrIdentity.Driver.log.Infof("VultrIdentityServer.Probe called with request : %v", req)

	return &csi.ProbeResponse{
		Ready: &wrappers.BoolValue{Value: true},
	}, nil
}
