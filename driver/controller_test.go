package driver

import (
	"context"
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
		client:       client,
		isController: true,
		log:          log,
		region:       "1",
	}

	return NewVultrControllerServer(d)
}
func TestCreateVolume(t *testing.T) {
	controller := NewFakeVultrControllerServer("create volume")

	_, err := controller.CreateVolume(context.TODO(), &csi.CreateVolumeRequest{
		Name: "volume-test-name",
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
	controller := NewFakeVultrControllerServer("delete volume")

	volumeID := "123456"
	_, err := controller.DeleteVolume(context.Background(), &csi.DeleteVolumeRequest{
		VolumeId: volumeID,
	})

	if err != nil {
		t.Errorf("Expected no error, got error : %v", err)
	}
}

func TestPublishVolume(t *testing.T) {
	controller := NewFakeVultrControllerServer("delete volume")

	nodeId := "123456"
	volumeID := "342512"

	_, err := controller.ControllerPublishVolume(context.Background(), &csi.ControllerPublishVolumeRequest{
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
}

func TestUnPublishVolume(t *testing.T) {
	controller := NewFakeVultrControllerServer("delete volume")

	nodeId := "123456"
	volumeID := "342512"

	_, err := controller.ControllerUnpublishVolume(context.Background(), &csi.ControllerUnpublishVolumeRequest{
		NodeId:   nodeId,
		VolumeId: volumeID,
	})

	if err != nil {
		t.Errorf("Expected no error, got error : %v", err)
	}
}
