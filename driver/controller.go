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
	"strconv"
	"strings"
	"time"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/sirupsen/logrus"
	"github.com/vultr/govultr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	_   = iota
	kiB = 1 << (10 * iota)
	miB
	giB
	tiB
)

const (
	minVolumeSizeInBytes      int64 = 1 * giB
	maxVolumeSizeInBytes      int64 = 10 * tiB
	defaultVolumeSizeInBytes  int64 = 10 * giB
	volumeStatusCheckRetries        = 10
	volumeStatusCheckInterval       = 1
)

var (
	supportedVolCapabilities = &csi.VolumeCapability_AccessMode{
		Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
	}
)

var _ csi.ControllerServer = &VultrControllerServer{}

type VultrControllerServer struct {
	Driver *VultrDriver
}

func NewVultrControllerServer(driver *VultrDriver) *VultrControllerServer {
	return &VultrControllerServer{Driver: driver}
}

// CreateVolume provisions a new volume on behalf of the user
func (c *VultrControllerServer) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	volName := req.Name

	if volName == "" {
		return nil, status.Error(codes.InvalidArgument, "CreateVolume Name is missing")
	}

	if req.VolumeCapabilities == nil || len(req.VolumeCapabilities) == 0 {
		return nil, status.Error(codes.InvalidArgument, "CreateVolume Volume Capabilities is missing")
	}

	// Validate
	if !isValidCapability(req.VolumeCapabilities) {
		return nil, status.Errorf(codes.InvalidArgument, "CreateVolume Volume capability is not compatible: %v", req)
	}

	// check that the volume doesnt already exist
	volumes, err := c.Driver.client.BlockStorage.List(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var curVolume *govultr.BlockStorage
	for _, volume := range volumes {
		if volume.Label == volName {
			curVolume = &volume
		}
	}

	if curVolume != nil {
		return &csi.CreateVolumeResponse{
			Volume: &csi.Volume{
				VolumeId:      curVolume.BlockStorageID,
				CapacityBytes: int64(curVolume.SizeGB) * giB,
			},
		}, nil
	}

	// if applicable, create volume
	region, err := strconv.Atoi(c.Driver.region)
	if err != nil {
		return nil, status.Error(codes.Aborted, "region code must be an int")
	}
	size, err := getStorageBytes(req.CapacityRange)
	if err != nil {
		return nil, status.Errorf(codes.OutOfRange, "invalid volume capacity range: %v", err)
	}

	c.Driver.log.WithFields(logrus.Fields{"size": size})

	volume, err := c.Driver.client.BlockStorage.Create(ctx, region, int(size/giB), volName)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Check to see if volume is in active state
	volReady := false

	for i := 0; i < volumeStatusCheckRetries; i++ {
		time.Sleep(volumeStatusCheckInterval * time.Second)
		bs, err := c.Driver.client.BlockStorage.Get(ctx, volume.BlockStorageID)

		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		if bs.Status == "active" {
			volReady = true
			break
		}
	}

	if !volReady {
		return nil, status.Errorf(codes.Internal, "volume is not active after %v seconds", volumeStatusCheckRetries)
	}

	res := &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			VolumeId:      volume.BlockStorageID,
			CapacityBytes: size,
			AccessibleTopology: []*csi.Topology{
				{
					Segments: map[string]string{
						"region": c.Driver.region,
					},
				},
			},
		},
	}

	return res, nil
}

func (c *VultrControllerServer) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "DeleteVolume VolumeID is missing")
	}

	list, err := c.Driver.client.BlockStorage.List(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	exists := false
	for _, v := range list {
		if v.BlockStorageID == req.VolumeId {
			exists = true
			break
		}
	}

	if !exists {
		return &csi.DeleteVolumeResponse{}, nil
	}

	err = c.Driver.client.BlockStorage.Delete(ctx, req.VolumeId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot delete volume, %v", err.Error())
	}

	return &csi.DeleteVolumeResponse{}, nil
}

func (c *VultrControllerServer) ControllerPublishVolume(ctx context.Context, req *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "ControllerPublishVolume Volume ID is missing")
	}

	if req.NodeId == "" {
		return nil, status.Error(codes.InvalidArgument, "ControllerPublishVolume Node ID is missing")
	}

	if req.VolumeCapability == nil {
		return nil, status.Error(codes.InvalidArgument, "ControllerPublishVolume Volume ID is missing")
	}

	if req.Readonly {
		return nil, status.Error(codes.InvalidArgument, "ControllerPublishVolume read only is not currently supported")
	}

	volume, err := c.Driver.client.BlockStorage.Get(ctx, req.VolumeId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "cannot get volume: %v", err.Error())
	}

	_, err = c.Driver.client.Server.GetServer(ctx, req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "cannot get node: %v", err.Error())
	}

	// node is already attached, do nothing
	if volume.InstanceID == req.NodeId {
		return &csi.ControllerPublishVolumeResponse{}, nil
	}

	// assuming its attached & to the wrong node
	if volume.InstanceID != "" {
		return nil, status.Errorf(codes.FailedPrecondition, "cannot attach volume to node because it is already attached to a different node ID: %v", volume.InstanceID)
	}

	err = c.Driver.client.BlockStorage.Attach(ctx, req.VolumeId, req.NodeId)
	if err != nil {
		// Desired node could still be spinning up
		if strings.Contains(err.Error(), "Server is currently locked") {
			return nil, status.Errorf(codes.Aborted, "cannot attach volume to node: %v", err.Error())
		}

		if strings.Contains(err.Error(), "Block storage volume is already attached to a node") {
			return &csi.ControllerPublishVolumeResponse{}, nil
		}
	}

	attachReady := false
	for i := 0; i < volumeStatusCheckRetries; i++ {
		time.Sleep(volumeStatusCheckInterval * time.Second)
		bs, err := c.Driver.client.BlockStorage.Get(ctx, volume.BlockStorageID)

		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		if bs.InstanceID == req.NodeId {
			attachReady = true
			break
		}
	}

	if !attachReady {
		return nil, status.Errorf(codes.Internal, "volume is not attached to node after %v seconds", volumeStatusCheckRetries)
	}

	return &csi.ControllerPublishVolumeResponse{}, nil
}

func (c *VultrControllerServer) ControllerUnpublishVolume(ctx context.Context, req *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "ControllerUnpublishVolume Volume ID is missing")
	}

	if req.NodeId != "" {
		return nil, status.Error(codes.InvalidArgument, "ControllerUnpublishVolume Node ID is missing")
	}

	volume, err := c.Driver.client.BlockStorage.Get(ctx, req.VolumeId)
	if err != nil {
		return &csi.ControllerUnpublishVolumeResponse{}, nil
	}

	// node is already unattached, do nothing
	if volume.InstanceID == "" {
		return &csi.ControllerUnpublishVolumeResponse{}, nil
	}

	_, err = c.Driver.client.Server.GetServer(ctx, req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "cannot get node: %v", err.Error())
	}

	err = c.Driver.client.BlockStorage.Detach(ctx, req.VolumeId)
	if err != nil {
		if strings.Contains(err.Error(), "Block storage volume is not currently attached to a server") {
			return &csi.ControllerUnpublishVolumeResponse{}, nil
		}
		return nil, status.Errorf(codes.Internal, "cannot detach volume: %v", err.Error())
	}

	return &csi.ControllerUnpublishVolumeResponse{}, nil
}

func (c *VultrControllerServer) ValidateVolumeCapabilities(context.Context, *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	panic("implement me")
}

func (c *VultrControllerServer) ListVolumes(ctx context.Context, req *csi.ListVolumesRequest) (*csi.ListVolumesResponse, error) {
	list, err := c.Driver.client.BlockStorage.List(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ListVolumes cannot retrieve list of volumes. %v", err.Error())
	}

	var entries []*csi.ListVolumesResponse_Entry
	for _, v := range list {
		entries = append(entries, &csi.ListVolumesResponse_Entry{
			Volume: &csi.Volume{
				VolumeId:      v.BlockStorageID,
				CapacityBytes: int64(v.SizeGB) * giB,
			},
		})
	}

	res := &csi.ListVolumesResponse{
		Entries: entries,
	}

	return res, nil
}

func (c *VultrControllerServer) GetCapacity(context.Context, *csi.GetCapacityRequest) (*csi.GetCapacityResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// ControllerGetCapabilities get capabilities of the controller
func (c *VultrControllerServer) ControllerGetCapabilities(context.Context, *csi.ControllerGetCapabilitiesRequest) (*csi.ControllerGetCapabilitiesResponse, error) {
	capability := func(capability csi.ControllerServiceCapability_RPC_Type) *csi.ControllerServiceCapability {
		return &csi.ControllerServiceCapability{
			Type: &csi.ControllerServiceCapability_Rpc{
				Rpc: &csi.ControllerServiceCapability_RPC{
					Type: capability,
				},
			},
		}
	}

	var capabilities []*csi.ControllerServiceCapability
	for _, caps := range []csi.ControllerServiceCapability_RPC_Type{
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
		csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME,
		csi.ControllerServiceCapability_RPC_LIST_VOLUMES,
	} {
		capabilities = append(capabilities, capability(caps))
	}

	resp := &csi.ControllerGetCapabilitiesResponse{
		Capabilities: capabilities,
	}

	c.Driver.log.WithFields(logrus.Fields{
		"response": resp,
		"method":   "controller-get-capabilities",
	})

	return resp, nil
}

func (c *VultrControllerServer) CreateSnapshot(context.Context, *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (c *VultrControllerServer) DeleteSnapshot(context.Context, *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (c *VultrControllerServer) ListSnapshots(context.Context, *csi.ListSnapshotsRequest) (*csi.ListSnapshotsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (c *VultrControllerServer) ControllerExpandVolume(context.Context, *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func isValidCapability(caps []*csi.VolumeCapability) bool {
	for _, capacity := range caps {
		if capacity == nil {
			return false
		}

		accessMode := capacity.GetAccessMode()
		if accessMode == nil {
			return false
		}

		if accessMode.GetMode() != supportedVolCapabilities.GetMode() {
			return false
		}

		accessType := capacity.GetAccessType()
		switch accessType.(type) {
		case *csi.VolumeCapability_Block:
		case *csi.VolumeCapability_Mount:
		default:
			return false
		}
	}
	return true
}

// getStorageBytes returns storage size in bytes
func getStorageBytes(capRange *csi.CapacityRange) (int64, error) {
	if capRange == nil {
		return defaultVolumeSizeInBytes, nil
	}

	capacity := capRange.GetRequiredBytes()
	return capacity, nil
}
