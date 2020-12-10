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
	"github.com/vultr/govultr/v2"
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
	volumeStatusCheckRetries        = 15
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

	c.Driver.log.WithFields(logrus.Fields{
		"volume-name":  volName,
		"capabilities": req.VolumeCapabilities,
	}).Info("Create Volume: called")

	// check that the volume doesnt already exist
	listOptions := &govultr.ListOptions{}
	var curVolume *govultr.BlockStorage

	for {
		volumes, meta, err := c.Driver.client.BlockStorage.List(ctx, listOptions)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		for _, v := range volumes {
			if v.Label == volName {
				curVolume = &v
				break
			}
		}

		if curVolume != nil {
			return &csi.CreateVolumeResponse{
				Volume: &csi.Volume{
					VolumeId:      curVolume.ID,
					CapacityBytes: int64(curVolume.SizeGB) * giB,
				},
			}, nil
		}

		if meta.Links.Next != "" {
			listOptions.Cursor = meta.Links.Next
			continue
		}

		break
	}

	// if applicable, create volume
	size, err := getStorageBytes(req.CapacityRange)
	if err != nil {
		return nil, status.Errorf(codes.OutOfRange, "invalid volume capacity range: %v", err)
	}

	blockReq := &govultr.BlockStorageCreate{
		Region: c.Driver.region,
		SizeGB: int(size / giB),
		Label:  volName,
	}

	volume, err := c.Driver.client.BlockStorage.Create(ctx, blockReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Check to see if volume is in active state
	volReady := false

	for i := 0; i < volumeStatusCheckRetries; i++ {
		time.Sleep(volumeStatusCheckInterval * time.Second)
		bs, err := c.Driver.client.BlockStorage.Get(ctx, volume.ID)

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
			VolumeId:      volume.ID,
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

	c.Driver.log.WithFields(logrus.Fields{
		"size":        size,
		"volume-id":   volume.ID,
		"volume-name": volume.Label,
		"volume-size": volume.SizeGB,
	}).Info("Create Volume: created volume")

	return res, nil
}

func (c *VultrControllerServer) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "DeleteVolume VolumeID is missing")
	}

	c.Driver.log.WithFields(logrus.Fields{
		"volume-id": req.VolumeId,
	}).Info("Delete volume: called")

	listOptions := &govultr.ListOptions{}
	exists := false
	for {
		list, meta, err := c.Driver.client.BlockStorage.List(ctx, listOptions)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		for _, v := range list {
			if v.ID == req.VolumeId {
				exists = true
				break
			}
		}

		if exists == true {
			break
		}

		if meta.Links.Next != "" {
			listOptions.Cursor = meta.Links.Next
			continue
		}
		c.Driver.log.WithFields(logrus.Fields{
			"volume-id": req.VolumeId,
		}).Info("!!!!!!!!!!!!!!!!!!!!!!!!!!")
		return &csi.DeleteVolumeResponse{}, nil
	}

	// detach just to be safe
	detach := &govultr.BlockStorageDetach{
		Live: govultr.BoolToBoolPtr(true),
	}
	err := c.Driver.client.BlockStorage.Detach(ctx, req.VolumeId, detach)
	if err != nil {
		if !strings.Contains(err.Error(), "Block storage volume is not currently attached to a server") {
			return nil, status.Errorf(codes.Internal, "cannot detach volume in delete, %v", err.Error())
		}
	}

	err = c.Driver.client.BlockStorage.Delete(ctx, req.VolumeId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot delete volume, %v", err.Error())
	}

	c.Driver.log.WithFields(logrus.Fields{
		"volume-id": req.VolumeId,
	}).Info("Delete Volume: deleted")

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
		return nil, status.Error(codes.InvalidArgument, "ControllerPublishVolume VolumeCapability is missing")
	}

	if req.Readonly {
		return nil, status.Error(codes.InvalidArgument, "ControllerPublishVolume read only is not currently supported")
	}

	volume, err := c.Driver.client.BlockStorage.Get(ctx, req.VolumeId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "cannot get volume: %v", err.Error())
	}

	_, err = c.Driver.client.Instance.Get(ctx, req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "cannot get node: %v", err.Error())
	}

	// node is already attached, do nothing
	if volume.AttachedToInstance == req.NodeId {
		return &csi.ControllerPublishVolumeResponse{
			PublishContext: map[string]string{
				c.Driver.publishVolumeID: volume.ID,
			},
		}, nil
	}

	// assuming its attached & to the wrong node
	if volume.AttachedToInstance != "" {
		return nil, status.Errorf(codes.FailedPrecondition, "cannot attach volume to node because it is already attached to a different node ID: %v", volume.AttachedToInstance)
	}

	c.Driver.log.WithFields(logrus.Fields{
		"volume-id": req.VolumeId,
		"node-id":   req.NodeId,
	}).Info("Controller Publish Volume: called")

	attach := &govultr.BlockStorageAttach{
		InstanceID: req.NodeId,
		Live:       govultr.BoolToBoolPtr(true),
	}
	err = c.Driver.client.BlockStorage.Attach(ctx, req.VolumeId, attach)
	if err != nil {
		// Desired node could still be spinning up
		if strings.Contains(err.Error(), "Server is currently locked") {
			return nil, status.Errorf(codes.Aborted, "cannot attach volume to node: %v", err.Error())
		}

		if strings.Contains(err.Error(), "Block storage volume is already attached to a server") {
			return &csi.ControllerPublishVolumeResponse{
				PublishContext: map[string]string{
					c.Driver.publishVolumeID: volume.ID,
				},
			}, nil
		}
	}

	attachReady := false
	for i := 0; i < volumeStatusCheckRetries; i++ {
		time.Sleep(volumeStatusCheckInterval * time.Second)
		bs, err := c.Driver.client.BlockStorage.Get(ctx, volume.ID)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		if bs.AttachedToInstance == req.NodeId {
			attachReady = true
			break
		}
	}

	if !attachReady {
		return nil, status.Errorf(codes.Internal, "volume is not attached to node after %v seconds", volumeStatusCheckRetries)
	}

	c.Driver.log.WithFields(logrus.Fields{
		"volume-id": req.VolumeId,
		"node-id":   req.NodeId,
	}).Info("Controller Publish Volume: published")

	return &csi.ControllerPublishVolumeResponse{
		PublishContext: map[string]string{
			c.Driver.publishVolumeID: volume.ID,
		},
	}, nil
}

func (c *VultrControllerServer) ControllerUnpublishVolume(ctx context.Context, req *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "ControllerUnpublishVolume Volume ID is missing")
	}

	if req.NodeId == "" {
		return nil, status.Error(codes.InvalidArgument, "ControllerUnpublishVolume Node ID is missing")
	}

	c.Driver.log.WithFields(logrus.Fields{
		"volume-id": req.VolumeId,
		"node-id":   req.NodeId,
	}).Info("Controller Publish Unpublish: called")

	volume, err := c.Driver.client.BlockStorage.Get(ctx, req.VolumeId)
	if err != nil {
		return &csi.ControllerUnpublishVolumeResponse{}, nil
	}

	// node is already unattached, do nothing
	if volume.AttachedToInstance == "" {
		return &csi.ControllerUnpublishVolumeResponse{}, nil
	}

	_, err = c.Driver.client.Instance.Get(ctx, req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "cannot get node: %v", err.Error())
	}
	detach := &govultr.BlockStorageDetach{
		Live: govultr.BoolToBoolPtr(true),
	}

	err = c.Driver.client.BlockStorage.Detach(ctx, req.VolumeId, detach)
	if err != nil {
		if strings.Contains(err.Error(), "Block storage volume is not currently attached to a server") {
			return &csi.ControllerUnpublishVolumeResponse{}, nil
		}
		return nil, status.Errorf(codes.Internal, "cannot detach volume: %v", err.Error())
	}

	c.Driver.log.WithFields(logrus.Fields{
		"volume-id": req.VolumeId,
		"node-id":   req.NodeId,
	}).Info("Controller Unublish Volume: unpublished")

	return &csi.ControllerUnpublishVolumeResponse{}, nil
}

// ValidateVolumeCapabilities checks if requested capabilities are supported
func (c *VultrControllerServer) ValidateVolumeCapabilities(ctx context.Context, req *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "ValidateVolumeCapabilities Volume ID is missing")
	}

	if req.VolumeCapabilities == nil {
		return nil, status.Error(codes.InvalidArgument, "ValidateVolumeCapabilities Volume Capabilities is missing")
	}

	_, err := c.Driver.client.BlockStorage.Get(ctx, req.VolumeId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "cannot get volume: %v", err.Error())
	}

	res := &csi.ValidateVolumeCapabilitiesResponse{
		Confirmed: &csi.ValidateVolumeCapabilitiesResponse_Confirmed{
			VolumeCapabilities: []*csi.VolumeCapability{
				{
					AccessMode: supportedVolCapabilities,
				},
			},
		},
	}

	return res, nil
}

func (c *VultrControllerServer) ListVolumes(ctx context.Context, req *csi.ListVolumesRequest) (*csi.ListVolumesResponse, error) {
	//todo setup paging
	if req.StartingToken != "" {
		_, err := strconv.Atoi(req.StartingToken)
		if err != nil {
			return nil, status.Errorf(codes.Aborted, "ListVolumes starting_token is invalid: %s", err)
		}
	}

	listOptions := &govultr.ListOptions{}
	var entries []*csi.ListVolumesResponse_Entry

	for {
		list, meta, err := c.Driver.client.BlockStorage.List(ctx, listOptions) //todo
		if err != nil {
			return nil, status.Errorf(codes.Internal, "ListVolumes cannot retrieve list of volumes. %v", err.Error())
		}
		for _, v := range list {
			entries = append(entries, &csi.ListVolumesResponse_Entry{
				Volume: &csi.Volume{
					VolumeId:      v.ID,
					CapacityBytes: int64(v.SizeGB) * giB,
				},
			})
		}

		if meta.Links.Next != "" {
			listOptions.Cursor = meta.Links.Next
			continue
		}
		break
	}

	res := &csi.ListVolumesResponse{
		Entries: entries,
	}

	c.Driver.log.WithFields(logrus.Fields{
		"volumes": entries,
	}).Info("List Volumes")
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
