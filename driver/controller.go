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
	"strconv"

	"github.com/container-storage-interface/spec/lib/go/csi"
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
	minVolumeSizeInBytes     int64 = 1 * giB
	maxVolumeSizeInBytes     int64 = 10 * tiB
	defaultVolumeSizeInBytes int64 = 10 * giB
)

var (
	supportedVolCapabilities = &csi.VolumeCapability_AccessMode{
		Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
	}
)

// CreateVolume provisions a new volume on behalf of the user
func (d *VultrDriver) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
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
	curVolume, err := d.client.BlockStorage.Get(context.TODO(), volName)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// if volume exists, do nothing - idempotency
	blockID, _ := strconv.Atoi(curVolume.BlockStorageID)
	if blockID != 0 {

		return &csi.CreateVolumeResponse{
			Volume: &csi.Volume{
				VolumeId:      curVolume.BlockStorageID,
				CapacityBytes: int64(curVolume.SizeGB) * giB,
			},
		}, nil
	}

	// if applicable, create volume
	label := "CSI Volume"
	region, _ := strconv.Atoi(d.region)
	size, err := getStorageBytes(req.CapacityRange)
	if err != nil {
		return nil, status.Errorf(codes.OutOfRange, "invalid volume capacity range: %v", err)
	}

	volume, err := d.client.BlockStorage.Create(ctx, region, int(size), label)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	res := &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			VolumeId:      volume.BlockStorageID,
			CapacityBytes: size,
			AccessibleTopology: []*csi.Topology{
				{
					Segments: map[string]string{
						"region": d.region,
					},
				},
			},
		},
	}

	return res, nil
}

// DeleteVolume deletes a volume created by CreateVolume
func DeleteVolume() {
	fmt.Print("IMPLEMENT ME")
}

// ControllerPublishVolume makes a volume available on a specified node. This will attach a volume to the node
func ControllerPublishVolume() {
	fmt.Print("IMPLEMENT ME")
}

// ControllerUnpublishVolume makes the volume unavailable on a given node. Makes a call to detach a volume from a node
func ControllerUnpublishVolume() {
	fmt.Print("IMPLEMENT ME")
}

// ValidateVolumeCapabilities returns capabilities of a volume
func ValidateVolumeCapabilities() {
	fmt.Print("IMPLEMENT ME")
}

// ListVolumes returns all available volumes
func ListVolumes() {
	fmt.Print("IMPLEMENT ME")
}

// GetCapacity returns capacity of total storage pool available
func GetCapacity() {
	fmt.Print("IMPLEMENT ME")
}

// ControllerGetCapabilities returns the capabilities of the Controller plugin
func ControllerGetCapabilities() {
	fmt.Print("IMPLEMENT ME")
}

// CreateSnapshot creates a snapshot
func CreateSnapshot() {
	fmt.Print("IMPLEMENT ME")
}

// DeleteSnapshot deletes a snapshot
func DeleteSnapshot() {
	fmt.Print("IMPLEMENT ME")
}

// ListSnapshots lists snapshots
func ListSnapshots() {
	fmt.Print("IMPLEMENT ME")
}

// ControllerExpandVolume ...
func ControllerExpandVolume() {
	fmt.Print("IMPLEMENT ME")
}

func isValidCapability(caps []*csi.VolumeCapability) bool {
	for _, cap := range caps {
		if cap == nil {
			return false
		}

		accessMode := cap.GetAccessMode()
		if accessMode == nil {
			return false
		}

		if accessMode.GetMode() != supportedVolCapabilities.GetMode() {
			return false
		}

		accessType := cap.GetAccessType()
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

	cap := capRange.GetRequiredBytes()
	return cap, nil
}
