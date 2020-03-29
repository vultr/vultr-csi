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
