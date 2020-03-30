package driver

import (
	"context"
	"testing"

	"github.com/container-storage-interface/spec/lib/go/csi"
)

func TestCreateVolume(t *testing.T) {
	d := &fakeStorageDriver{}

	_, err := d.CreateVolume(context.Background(), &csi.CreateVolumeRequest{
		Name: "volume-name",
		VolumeCapabilities: []*csi.VolumeCapability{
			&csi.VolumeCapability{
				AccessType: &csi.VolumeCapability_Mount{
					Mount: &csi.VolumeCapability_MountVolume{},
				},
				AccessMode: &csi.VolumeCapability_AccessMode{
					Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
				},
			},
		},
	})

	if err != nil {
		t.Errorf("Expected no error, got error : %v", err)
	}
}

func TestDeleteVolume(t *testing.T) {
	d := &fakeStorageDriver{}

	volumeID := "123456"
	_, err := d.DeleteVolume(context.Background(), &csi.DeleteVolumeRequest{
		VolumeId: volumeID,
	})

	if err != nil {
		t.Errorf("Expected no error, got error : %v", err)
	}
}

func TestPublishVolume(t *testing.T) {
	d := &fakeStorageDriver{}
	nodeId := "111111"
	volumeID := "123456"

	_, err := d.ControllerPublishVolume(context.Background(), &csi.ControllerPublishVolumeRequest{
		NodeId:   nodeId,
		VolumeId: volumeID,
	})

	if err != nil {
		t.Errorf("Expected no error, got error : %v", err)
	}
}

func TestUnPublishVolume(t *testing.T) {
	d := &fakeStorageDriver{}
	nodeId := "111111"
	volumeID := "123456"

	_, err := d.ControllerUnPublishVolume(context.Background(), &csi.ControllerUnpublishVolumeRequest{
		NodeId:   nodeId,
		VolumeId: volumeID,
	})

	if err != nil {
		t.Errorf("Expected no error, got error : %v", err)
	}
}
