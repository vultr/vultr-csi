package driver

import (
	"context"
	"reflect"
	"testing"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/sirupsen/logrus"
)

func NewFakeVultrControllerServer(testName string) *VultrControllerServer {
	client := newFakeClient()
	log := logrus.New().WithFields(logrus.Fields{
		"test": testName,
	})

	d := &VultrDriver{
		client:          client,
		isController:    true,
		log:             log,
		region:          "ewr",
		publishVolumeID: "c56c7b6e-15c2-445e-9a5d-1063ab5828ec",
	}

	return NewVultrControllerServer(d)
}

func TestControllerCreateBlockVolume(t *testing.T) {
	controller := NewFakeVultrControllerServer("create block volume")

	res, err := controller.CreateVolume(context.TODO(), &csi.CreateVolumeRequest{
		Name: "volume-test-name",
		Parameters: map[string]string{
			"storage_type": "block",
			"disk_type":    "hdd",
		},
		VolumeCapabilities: []*csi.VolumeCapability{
			{
				AccessMode: &csi.VolumeCapability_AccessMode{
					Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
				},
			},
		},
		// CapacityRange: &csi.CapacityRange{
		// 	RequiredBytes: 42949672960,
		// },
	})

	if err != nil {
		t.Errorf("got error, expected no error: %v", err)
	}

	expected := &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			VolumeId:      "a35badcb-a4db-4171-9b9a-11910dfdb8f3",
			CapacityBytes: 42949672960,
			AccessibleTopology: []*csi.Topology{
				{
					Segments: map[string]string{
						"region": "ewr",
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(res, expected) {
		t.Errorf("expected %+v got %+v", expected, res)
	}
}

func TestControllerDeleteBlockVolume(t *testing.T) {
	controller := NewFakeVultrControllerServer("delete block volume")

	volumeID := "c56c7b6e-15c2-445e-9a5d-1063ab5828ec" //nolint:goconst
	res, err := controller.DeleteVolume(context.Background(), &csi.DeleteVolumeRequest{
		VolumeId: volumeID,
	})

	if err != nil {
		t.Errorf("Expected no error, got error : %v", err)
	}

	expected := &csi.DeleteVolumeResponse{}

	if !reflect.DeepEqual(res, expected) {
		t.Errorf("expected %+v got %+v", res, expected)
	}
}

func TestControllerPublishBlockVolume(t *testing.T) {
	controller := NewFakeVultrControllerServer("publish block volume")

	nodeID := "245bb2fe-b55c-44a0-9a1e-ab80e4b5f088" //nolint:goconst
	volumeID := "c56c7b6e-15c2-445e-9a5d-1063ab5828ec"

	res, err := controller.ControllerPublishVolume(context.Background(), &csi.ControllerPublishVolumeRequest{
		NodeId:   nodeID,
		VolumeId: volumeID,
		VolumeCapability: &csi.VolumeCapability{
			AccessType: &csi.VolumeCapability_Mount{
				Mount: &csi.VolumeCapability_MountVolume{},
			},
			AccessMode: &csi.VolumeCapability_AccessMode{
				Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
			},
		},
	})

	if err != nil {
		t.Errorf("Expected no error, got error : %v", err)
	}

	expected := &csi.ControllerPublishVolumeResponse{
		PublishContext: map[string]string{
			"mount_vol_name": "test-mount-3",
			"storage_type":   "block",
		},
	}

	if !reflect.DeepEqual(res, expected) {
		t.Errorf("expected %+v got %+v", res, expected)
	}
}

func TestControllerUnpublishBlockVolume(t *testing.T) {
	controller := NewFakeVultrControllerServer("unpublish block volume")

	nodeID := "c56c7b6e-15c2-445e-9a5d-1063ab5828ec"
	volumeID := "245bb2fe-b55c-44a0-9a1e-ab80e4b5f088"

	res, err := controller.ControllerUnpublishVolume(context.Background(), &csi.ControllerUnpublishVolumeRequest{
		NodeId:   nodeID,
		VolumeId: volumeID,
	})

	if err != nil {
		t.Errorf("Expected no error, got error : %v", err)
	}

	expected := &csi.ControllerUnpublishVolumeResponse{}

	if !reflect.DeepEqual(res, expected) {
		t.Errorf("expected %+v got %+v", res, expected)
	}
}

func TestControllerUnpublishBlockVolumeNotFound(t *testing.T) {
	controller := NewFakeVultrControllerServer("unpublish block volume not found")

	nodeID := "c56c7b6e-15c2-445e-9a5d-1063ab5828ec"
	// Use a volume ID that doesn't exist in any storage type
	volumeID := "nonexistent-volume-id-12345"

	res, err := controller.ControllerUnpublishVolume(context.Background(), &csi.ControllerUnpublishVolumeRequest{
		NodeId:   nodeID,
		VolumeId: volumeID,
	})

	// Should succeed even though volume doesn't exist (idempotency)
	if err != nil {
		t.Errorf("Expected no error when unpublishing non-existent volume, got error : %v", err)
	}

	expected := &csi.ControllerUnpublishVolumeResponse{}

	if !reflect.DeepEqual(res, expected) {
		t.Errorf("expected %+v got %+v", res, expected)
	}
}
