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
		region:          "1",
		publishVolumeID: "342512",
	}

	return NewVultrControllerServer(d)
}
func TestCreateVolume(t *testing.T) {
	controller := NewFakeVultrControllerServer("create volume")

	res, err := controller.CreateVolume(context.TODO(), &csi.CreateVolumeRequest{
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
		t.Errorf("got error, expected no error: %v", err)
	}

	expected := &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			VolumeId:      "342512",
			CapacityBytes: 10737418240,
			AccessibleTopology: []*csi.Topology{
				{
					Segments: map[string]string{
						"region": "1",
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(expected.Volume.AccessibleTopology, res.Volume.AccessibleTopology) {
		t.Errorf("expected %+v got %+v", res, expected)
	}

	if expected.Volume.CapacityBytes != res.Volume.CapacityBytes {
		t.Errorf("expected %+v got %+v", res, expected)
	}
}

func TestDeleteVolume(t *testing.T) {
	controller := NewFakeVultrControllerServer("delete volume")

	volumeID := "342512"
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

	create, _ := controller.CreateVolume(context.TODO(), &csi.CreateVolumeRequest{
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

	res, err := controller.ControllerPublishVolume(context.Background(), &csi.ControllerPublishVolumeRequest{
		NodeId:   "123456",
		VolumeId: create.Volume.VolumeId,
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
			controller.Driver.publishVolumeID: create.Volume.VolumeId,
		},
	}

	if !reflect.DeepEqual(res, expected) {
		t.Errorf("expected %+v got %+v", res, expected)
	}
}

func TestUnPublishVolume(t *testing.T) {
	controller := NewFakeVultrControllerServer("delete volume")

	nodeId := "123456"
	volumeID := "342512"

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
