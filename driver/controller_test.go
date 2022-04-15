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
func TestCreateVolume(t *testing.T) {
	controller := NewFakeVultrControllerServer("create volume")

	res, err := controller.CreateVolume(context.TODO(), &csi.CreateVolumeRequest{
		Name:       "volume-test-name",
		Parameters: map[string]string{"block_type": "high_perf"},
		VolumeCapabilities: []*csi.VolumeCapability{
			{
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
		t.Errorf("got error, expected no error: %v", err)
	}

	expected := &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			VolumeId:      "c56c7b6e-15c2-445e-9a5d-1063ab5828ec",
			CapacityBytes: 10737418240,
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
		t.Errorf("expected %+v got %+v", res, expected)
	}
}

func TestDeleteVolume(t *testing.T) {
	controller := NewFakeVultrControllerServer("delete volume")

	volumeID := "c56c7b6e-15c2-445e-9a5d-1063ab5828ec"
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

func TestPublishVolume(t *testing.T) {
	controller := NewFakeVultrControllerServer("delete volume")

	nodeId := "245bb2fe-b55c-44a0-9a1e-ab80e4b5f088"
	volumeID := "c56c7b6e-15c2-445e-9a5d-1063ab5828ec"

	res, err := controller.ControllerPublishVolume(context.Background(), &csi.ControllerPublishVolumeRequest{
		NodeId:   nodeId,
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
			controller.Driver.publishVolumeID: volumeID,
		},
	}

	if !reflect.DeepEqual(res, expected) {
		t.Errorf("expected %+v got %+v", res, expected)
	}
}

func TestUnPublishVolume(t *testing.T) {
	controller := NewFakeVultrControllerServer("delete volume")

	nodeId := "c56c7b6e-15c2-445e-9a5d-1063ab5828ec"
	volumeID := "245bb2fe-b55c-44a0-9a1e-ab80e4b5f088"

	res, err := controller.ControllerUnpublishVolume(context.Background(), &csi.ControllerUnpublishVolumeRequest{
		NodeId:   nodeId,
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
