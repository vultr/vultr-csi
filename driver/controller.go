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
	"fmt"
	"strings"
	"time"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/sirupsen/logrus"
	"github.com/vultr/vultr-csi/internal/vultrstorage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	gibiByte                  int64 = 1073741824
	volumeStatusCheckRetries  int   = 15
	volumeStatusCheckInterval int   = 1
)

var _ csi.ControllerServer = &VultrControllerServer{}

// VultrControllerServer is the struct type for the VultrDriver
type VultrControllerServer struct {
	csi.UnimplementedControllerServer
	Driver *VultrDriver
}

// NewVultrControllerServer returns a VultrControllerServer
func NewVultrControllerServer(driver *VultrDriver) *VultrControllerServer {
	return &VultrControllerServer{Driver: driver}
}

// CreateVolume provisions a new volume on behalf of the user
func (c *VultrControllerServer) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) { //nolint:gocyclo,lll,funlen
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "CreateVolume: name is missing")
	}

	if len(req.VolumeCapabilities) == 0 {
		return nil, status.Error(codes.InvalidArgument, "CreateVolume: capabilities is missing")
	}

	diskType := req.Parameters["disk_type"]
	storageType := req.Parameters["storage_type"]
	blockType := req.Parameters["block_type"]

	if diskType == "" {
		return nil, status.Error(codes.InvalidArgument, "CreateVolume: parameter `disk_type` is missing")
	}

	if storageType == "" {
		return nil, status.Error(codes.InvalidArgument, "CreateVolume: parameter `storage_type` is missing")
	}

	// handle legacy param
	if blockType != "" {
		diskType = "block"
		if blockType == "high_perf" {
			storageType = "nvme"
		} else if blockType == "storage_opt" {
			storageType = "hdd"
		}
	}

	sh, err := vultrstorage.NewVultrStorageHandler(c.Driver.client, storageType, diskType)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "CreateVolume: cannot initialize vultr storage handler: %v", err.Error())
	}

	if err := validateCapabilities(req.VolumeCapabilities, sh.Capabilities); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "CreateVolume: requested capability is not compatible: %v", err)
	}

	c.Driver.log.WithFields(logrus.Fields{
		"volume-name":  req.Name,
		"capabilities": req.VolumeCapabilities,
	}).Info("CreateVolume: called")

	var curVolume *vultrstorage.VultrStorage

	storages, err := vultrstorage.ListAllStorages(ctx, c.Driver.client)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "CreateVolume: could not retrieve list of storages. %v", err.Error())
	}

	// check if volume already exists
	for i := range storages {
		if storages[i].Label == req.Name {
			curVolume = &storages[i]
			break
		}
	}

	if curVolume != nil {
		return &csi.CreateVolumeResponse{
			Volume: &csi.Volume{
				VolumeId:      curVolume.ID,
				CapacityBytes: int64(curVolume.SizeGB) * gibiByte,
			},
		}, nil
	}

	// volume doesn't exist, create
	size, err := getStorageBytes(req.CapacityRange, sh)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "CreateVolume: could not request new volume: %v", err.Error())
	}

	storageReq := &vultrstorage.VultrStorageReq{
		Region:   c.Driver.region,
		SizeGB:   int(size / gibiByte),
		Label:    req.Name,
		DiskType: diskType,
	}

	volume, err := sh.Operations.Create(ctx, *storageReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "CreateVolume: could not create a new volume: %v", err.Error())
	}

	// Check to see if volume is in active state
	volReady := false

	for i := 0; i < volumeStatusCheckRetries; i++ {
		time.Sleep(time.Duration(volumeStatusCheckInterval) * time.Second)

		storage, err := sh.Operations.Get(ctx, volume.ID)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "CreateVolume: could not retrieve the new volume: %v", err.Error())
		}

		if storage.Status == "active" {
			volReady = true
			break
		}
	}

	if !volReady {
		return nil, status.Errorf(codes.Internal, "CreateVolume: volume is not active after %v seconds", volumeStatusCheckRetries)
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
	}).Info("CreateVolume: created volume")

	return res, nil
}

// DeleteVolume performs the volume deletion
func (c *VultrControllerServer) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "DeleteVolume: volume ID is missing")
	}

	c.Driver.log.WithFields(logrus.Fields{
		"volume-id": req.VolumeId,
	}).Info("DeleteVolume: called")

	exists := false
	var deleteStorage vultrstorage.VultrStorage

	storages, err := vultrstorage.ListAllStorages(ctx, c.Driver.client)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DeleteVolume: could not retrieve list of storages. %v", err.Error())
	}

	for i := range storages {
		if storages[i].ID == req.VolumeId {
			exists = true
			deleteStorage = storages[i]
			break
		}
	}

	if !exists {
		return &csi.DeleteVolumeResponse{}, nil
	}

	sh, err := vultrstorage.NewVultrStorageHandler(c.Driver.client, deleteStorage.StorageType, "")
	if err != nil {
		return nil, fmt.Errorf("DeleteVolume: cannot initialize vultr storage handler. %v", err)
	}

	// detach all instances
	for i := range deleteStorage.AttachedInstances {
		if err := sh.Operations.Detach(ctx, deleteStorage.ID, deleteStorage.AttachedInstances[i].NodeID); err != nil {
			if !strings.Contains(err.Error(), "volume is not currently attached") ||
				!strings.Contains(err.Error(), "Attachment Not Found") {
				return nil, status.Errorf(codes.Internal, "DeleteVolume: cannot detach volume in delete, %v", err.Error())
			}
		}
	}

	// otherwise, internal brokenness
	if err := sh.Operations.Delete(ctx, deleteStorage.ID); err != nil {
		return nil, status.Errorf(codes.Internal, "DeleteVolume: cannot delete volume, %v", err.Error())
	}

	c.Driver.log.WithFields(logrus.Fields{
		"volume-id": req.VolumeId,
	}).Info("DeleteVolume: deleted")

	return &csi.DeleteVolumeResponse{}, nil
}

// ControllerPublishVolume performs the volume publish for the controller
func (c *VultrControllerServer) ControllerPublishVolume(ctx context.Context, req *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) { //nolint:lll,gocyclo
	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "ControllerPublishVolume: volume ID is missing")
	}

	if req.NodeId == "" {
		return nil, status.Error(codes.InvalidArgument, "ControllerPublishVolume: node ID is missing")
	}

	if req.VolumeCapability == nil {
		return nil, status.Error(codes.InvalidArgument, "ControllerPublishVolume: volume capability is missing")
	}

	if req.Readonly {
		return nil, status.Error(codes.InvalidArgument, "ControllerPublishVolume: read only is not currently supported")
	}

	sh, err := vultrstorage.FindVultrStorageHandlerByID(ctx, c.Driver.client, req.VolumeId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ControllerPublishVolume: could not find storage handler for storage. %v", err.Error())
	}

	storageExisting, err := sh.Operations.Get(ctx, req.VolumeId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "ControllerPublishVolume: could not retrieve existing storage volume: %v", err.Error())
	}

	if _, _, err = c.Driver.client.Instance.Get(ctx, req.NodeId); err != nil { //nolint:bodyclose
		return nil, status.Errorf(codes.NotFound, "ControllerPublishVolume: could not retrieve node: %v", err.Error())
	}

	for i := range storageExisting.AttachedInstances {
		if storageExisting.AttachedInstances[i].NodeID == req.NodeId {
			return &csi.ControllerPublishVolumeResponse{
				PublishContext: map[string]string{
					"mount_vol_name": storageExisting.AttachedInstances[i].MountName,
					"storage_type":   storageExisting.StorageType,
				},
			}, nil
		}
	}

	// block storage cannot be mounted to more than one instance
	if storageExisting.StorageType == "block" && len(storageExisting.AttachedInstances) > 0 {
		return nil, status.Errorf(codes.FailedPrecondition,
			"ControllerPublishVolume: cannot attach volume to node because it is already attached to a different node ID: %v",
			storageExisting.AttachedInstances[0].NodeID)
	}

	c.Driver.log.WithFields(logrus.Fields{
		"volume-id": req.VolumeId,
		"node-id":   req.NodeId,
	}).Info("ControllerPublishVolume: called")

	err = sh.Operations.Attach(ctx, req.VolumeId, req.NodeId)
	if err != nil {
		if strings.Contains(err.Error(), "Server is currently locked") {
			return nil, status.Errorf(codes.Aborted, "cannot attach volume to node: %v", err.Error())
		}

		return nil, status.Errorf(codes.Internal, "ControllPublishVolume: cannot attach volume to node: %v", err.Error())
	}

	attachReady := false
	var storageAttached *vultrstorage.VultrStorage
	publishedVolName := ""

retries:
	for i := 0; i < volumeStatusCheckRetries; i++ {
		time.Sleep(time.Duration(volumeStatusCheckInterval) * time.Second)
		storageAttached, err = sh.Operations.Get(ctx, storageExisting.ID)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				"ControllerPublishVolume: unable to retrieve storage for retry check: %v",
				err.Error(),
			)
		}

		for i := range storageAttached.AttachedInstances {
			if storageAttached.AttachedInstances[i].NodeID == req.NodeId {
				attachReady = true
				publishedVolName = storageAttached.AttachedInstances[i].MountName
				break retries
			}
		}
	}

	if !attachReady {
		return nil, status.Errorf(
			codes.Internal,
			"ControllerPublishVolume: volume is not attached to node after %v seconds",
			volumeStatusCheckRetries,
		)
	}

	c.Driver.log.WithFields(logrus.Fields{
		"volume-id": req.VolumeId,
		"node-id":   req.NodeId,
	}).Info("ControllerPublishVolume: published")

	return &csi.ControllerPublishVolumeResponse{
		PublishContext: map[string]string{
			"mount_vol_name": publishedVolName,
			"storage_type":   storageAttached.StorageType,
		},
	}, nil
}

// ControllerUnpublishVolume performs the volume un-publish
func (c *VultrControllerServer) ControllerUnpublishVolume(ctx context.Context, req *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) { //nolint:lll
	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "ControllerUnpublishVolume: volume ID is missing")
	}

	if req.NodeId == "" {
		return nil, status.Error(codes.InvalidArgument, "ControllerUnpublishVolume: node ID is missing")
	}

	c.Driver.log.WithFields(logrus.Fields{
		"volume-id": req.VolumeId,
		"node-id":   req.NodeId,
	}).Info("ControllerPublishUnpublish: called")

	sh, err := vultrstorage.FindVultrStorageHandlerByID(ctx, c.Driver.client, req.VolumeId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ControllerUnpublishVolume: could not find storage handler for storage. %v", err.Error())
	}

	storage, err := sh.Operations.Get(ctx, req.VolumeId)
	if err != nil {
		// Not found, return empty response
		return &csi.ControllerUnpublishVolumeResponse{}, nil
	}

	// node is already unattached, do nothing
	if len(storage.AttachedInstances) == 0 {
		return &csi.ControllerUnpublishVolumeResponse{}, nil
	}

	for i := range storage.AttachedInstances {
		if storage.AttachedInstances[i].NodeID != req.NodeId {
			continue
		}

		if err := sh.Operations.Detach(ctx, req.VolumeId, storage.AttachedInstances[i].NodeID); err != nil {
			if strings.Contains(err.Error(), "Block storage volume is not currently attached to a server") ||
				strings.Contains(err.Error(), "Attachment Not Found") {
				return &csi.ControllerUnpublishVolumeResponse{}, nil
			}

			return nil, status.Errorf(codes.Internal, "ControllerUnpublishVolume: could not detach volume: %v", err.Error())
		}
	}

	c.Driver.log.WithFields(logrus.Fields{
		"volume-id": req.VolumeId,
		"node-id":   req.NodeId,
	}).Info("ControllerUnublishVolume: unpublished")

	return &csi.ControllerUnpublishVolumeResponse{}, nil
}

// ControllerModifyVolume is unimplemented
func (c *VultrControllerServer) ControllerModifyVolume(ctx context.Context, req *csi.ControllerModifyVolumeRequest) (*csi.ControllerModifyVolumeResponse, error) { //nolint:lll
	return nil, status.Errorf(codes.Unimplemented, "method ControllerModifyVolume not implemented")
}

// ValidateVolumeCapabilities checks if requested capabilities are supported
func (c *VultrControllerServer) ValidateVolumeCapabilities(ctx context.Context, req *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) { //nolint:lll
	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "ValidateVolumeCapabilities: volume ID is missing")
	}

	if req.VolumeCapabilities == nil {
		return nil, status.Error(codes.InvalidArgument, "ValidateVolumeCapabilities: volume Capabilities is missing")
	}

	diskType := req.Parameters["disk_type"]
	storageType := req.Parameters["storage_type"]
	blockType := req.Parameters["block_type"]

	// handle legacy param
	if blockType != "" {
		diskType = "block"
		if blockType == "high_perf" {
			storageType = "nvme"
		} else if blockType == "storage_opt" {
			storageType = "hdd"
		}
	}

	sh, err := vultrstorage.NewVultrStorageHandler(c.Driver.client, storageType, diskType)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ValidateVolumeCapabilities: cannot initialize vultr storage handler. %v", err.Error())
	}

	if _, err := sh.Operations.Get(ctx, req.VolumeId); err != nil {
		return nil, status.Errorf(codes.NotFound, "cannot get volume: %v", err.Error())
	}

	return &csi.ValidateVolumeCapabilitiesResponse{
		Confirmed: &csi.ValidateVolumeCapabilitiesResponse_Confirmed{VolumeCapabilities: sh.Capabilities},
	}, nil
}

// ListVolumes performs the list volume function
func (c *VultrControllerServer) ListVolumes(ctx context.Context, req *csi.ListVolumesRequest) (*csi.ListVolumesResponse, error) {
	var entries []*csi.ListVolumesResponse_Entry

	storages, err := vultrstorage.ListAllStorages(ctx, c.Driver.client)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ListVolumes: cannot retrieve all volumes: %v", err.Error())
	}

	for i := range storages {
		entries = append(entries, &csi.ListVolumesResponse_Entry{
			Volume: &csi.Volume{
				VolumeId:      storages[i].ID,
				CapacityBytes: int64(storages[i].SizeGB) * gibiByte,
			},
		})
	}

	res := &csi.ListVolumesResponse{
		Entries: entries,
	}

	c.Driver.log.WithFields(logrus.Fields{
		"volumes": entries,
	}).Info("ListVolumes: called")

	return res, nil
}

func (c *VultrControllerServer) GetCapacity(context.Context, *csi.GetCapacityRequest) (*csi.GetCapacityResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// ControllerGetCapabilities get capabilities of the controller
func (c *VultrControllerServer) ControllerGetCapabilities(context.Context, *csi.ControllerGetCapabilitiesRequest) (*csi.ControllerGetCapabilitiesResponse, error) { //nolint:lll
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
		csi.ControllerServiceCapability_RPC_EXPAND_VOLUME,
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

// CreateSnapshot provides snapshot creation
func (c *VultrControllerServer) CreateSnapshot(context.Context, *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// DeleteSnapshot provides snapshot deletion
func (c *VultrControllerServer) DeleteSnapshot(context.Context, *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// ListSnapshots provides the list snapshot
func (c *VultrControllerServer) ListSnapshots(context.Context, *csi.ListSnapshotsRequest) (*csi.ListSnapshotsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// ControllerExpandVolume provides the expand volume
func (c *VultrControllerServer) ControllerExpandVolume(ctx context.Context, req *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) { //nolint:lll
	volumeID := req.GetVolumeId()
	if volumeID == "" {
		return nil, status.Error(codes.InvalidArgument, "ControllerExpandVolume: volume ID must be provided")
	}

	sh, err := vultrstorage.FindVultrStorageHandlerByID(ctx, c.Driver.client, req.VolumeId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ControllerExpandVolume: could not find storage handler for volume: %v", err.Error())
	}

	curVolume, err := sh.Operations.Get(ctx, req.VolumeId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ControllerExpandVolume: could not retrieve volume: %v", err.Error())
	}

	newSizeBytes, err := getStorageBytes(req.CapacityRange, sh)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ControllerExpandVolume: unable to determine requested size: %v", err.Error())
	}

	if int64(curVolume.SizeGB)*gibiByte > newSizeBytes {
		return nil, status.Error(codes.InvalidArgument, "ControllerExpandVolume: requested size must be larger than current size.")
	}

	c.Driver.log.WithFields(logrus.Fields{
		"volume-id": req.VolumeId,
		"size":      newSizeBytes,
	}).Info("ControllerExpandVolume: called")

	updateReq := vultrstorage.VultrStorageUpdateReq{
		SizeGB: int(newSizeBytes / gibiByte),
	}

	if _, err := sh.Operations.Update(ctx, req.VolumeId, updateReq); err != nil {
		return nil, status.Errorf(codes.Internal, "ControllerExpandVolume: unable to update storage: %v", err.Error())
	}

	nodeExpansion := false
	if sh.StorageType == "block" {
		nodeExpansion = true
	}

	return &csi.ControllerExpandVolumeResponse{CapacityBytes: newSizeBytes, NodeExpansionRequired: nodeExpansion}, nil
}

// ControllerGetVolume This relates to being able to get health checks on a PV. We do not have this
func (c *VultrControllerServer) ControllerGetVolume(ctx context.Context, request *csi.ControllerGetVolumeRequest) (*csi.ControllerGetVolumeResponse, error) { //nolint:lll
	return nil, status.Error(codes.Unimplemented, "")
}

// validateCapabilities compares the requested capabilities with the supported
// capabilities and returning false if not supported. Currently, only the
// AccessMode capability is relevant and checked.
func validateCapabilities(reqCaps, storageCaps []*csi.VolumeCapability) error {
	for i := range reqCaps {
		if reqCaps[i] == nil {
			return fmt.Errorf("requested capability missing")
		}

		accessMode := reqCaps[i].GetAccessMode()
		if accessMode == nil {
			return fmt.Errorf("requested capability access mode is missing")
		}

		var modeMatch bool
		for j := range storageCaps {
			if accessMode.GetMode() == storageCaps[j].AccessMode.GetMode() {
				modeMatch = true
			}
		}

		if !modeMatch {
			return fmt.Errorf("requested capability access mode is not supported")
		}

		accessType := reqCaps[i].GetAccessType()
		if accessType != nil {
			switch accessType.(type) {
			case *csi.VolumeCapability_Block:
			case *csi.VolumeCapability_Mount:
			default:
				return fmt.Errorf("requested capability is not supported: %v", accessType)
			}
		}
	}

	return nil
}

func getStorageBytes(capRange *csi.CapacityRange, sh *vultrstorage.VultrStorageHandler) (int64, error) {
	// return the csi capacity in bytes if present
	if capRange != nil {
		return capRange.GetRequiredBytes(), nil
	}

	// otherwise return the defaults
	if sh.DefaultSize != 0 {
		return sh.DefaultSize, nil
	}

	return 0, fmt.Errorf("default size unavailable for type %v storage %v disk", sh.StorageType, sh.DiskType)
}
