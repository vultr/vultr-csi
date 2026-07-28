package driver

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func TestControllerCreateBlockVolumeFromSnapshot(t *testing.T) {
	controller := NewFakeVultrControllerServer("create block volume from snapshot")

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
		VolumeContentSource: &csi.VolumeContentSource{
			Type: &csi.VolumeContentSource_Snapshot{
				Snapshot: &csi.VolumeContentSource_SnapshotSource{
					SnapshotId: "cb676a46-66fd-4dfb-b839-443f2e6c0b60",
				},
			},
		},
	})

	if err != nil {
		t.Errorf("got error, expected no error: %v", err)
	}

	expected := &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			VolumeId:      "a35badcb-a4db-4171-9b9a-11910dfdb8f3",
			CapacityBytes: 42949672960,
			ContentSource: &csi.VolumeContentSource{
				Type: &csi.VolumeContentSource_Snapshot{
					Snapshot: &csi.VolumeContentSource_SnapshotSource{
						SnapshotId: "cb676a46-66fd-4dfb-b839-443f2e6c0b60",
					},
				},
			},
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

func TestControllerCreateSnapshot(t *testing.T) {
	controller := NewFakeVultrControllerServer("create snapshot")

	res, err := controller.CreateSnapshot(context.Background(), &csi.CreateSnapshotRequest{
		Name:           "new-snapshot-test-name",
		SourceVolumeId: "c56c7b6e-15c2-445e-9a5d-1063ab5828ec",
	})

	if err != nil {
		t.Errorf("Expected no error, got error : %v", err)
	}

	createdAt, _ := time.Parse("2006-01-02 15:04:05", "2026-07-17 16:46:05")
	expected := &csi.CreateSnapshotResponse{
		Snapshot: &csi.Snapshot{
			SnapshotId:     "cb676a46-66fd-4dfb-b839-443f2e6c0b60",
			SourceVolumeId: "c56c7b6e-15c2-445e-9a5d-1063ab5828ec",
			SizeBytes:      10737418240,
			CreationTime:   timestamppb.New(createdAt),
			ReadyToUse:     true,
		},
	}

	if !reflect.DeepEqual(res, expected) {
		t.Errorf("expected %+v got %+v", expected, res)
	}
}

func TestControllerDeleteSnapshot(t *testing.T) {
	controller := NewFakeVultrControllerServer("delete snapshot")

	res, err := controller.DeleteSnapshot(context.Background(), &csi.DeleteSnapshotRequest{
		SnapshotId: "cb676a46-66fd-4dfb-b839-443f2e6c0b60",
	})

	if err != nil {
		t.Errorf("Expected no error, got error : %v", err)
	}

	expected := &csi.DeleteSnapshotResponse{}

	if !reflect.DeepEqual(res, expected) {
		t.Errorf("expected %+v got %+v", expected, res)
	}
}

func TestControllerListSnapshots(t *testing.T) {
	controller := NewFakeVultrControllerServer("list snapshots")

	res, err := controller.ListSnapshots(context.Background(), &csi.ListSnapshotsRequest{
		SourceVolumeId: "c56c7b6e-15c2-445e-9a5d-1063ab5828ec",
	})

	if err != nil {
		t.Errorf("Expected no error, got error : %v", err)
	}

	createdAt, _ := time.Parse("2006-01-02 15:04:05", "2026-07-17 16:46:05")
	expected := &csi.ListSnapshotsResponse{
		Entries: []*csi.ListSnapshotsResponse_Entry{
			{
				Snapshot: &csi.Snapshot{
					SnapshotId:     "cb676a46-66fd-4dfb-b839-443f2e6c0b60",
					SourceVolumeId: "c56c7b6e-15c2-445e-9a5d-1063ab5828ec",
					SizeBytes:      10737418240,
					CreationTime:   timestamppb.New(createdAt),
					ReadyToUse:     true,
				},
			},
		},
	}

	if !reflect.DeepEqual(res, expected) {
		t.Errorf("expected %+v got %+v", expected, res)
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
