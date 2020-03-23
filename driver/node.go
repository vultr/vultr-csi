package driver

import (
	"context"
	"fmt"
	csi "github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"path/filepath"
)

const (
	diskPath   = "/dev/disk/by-id"
	diskPrefix = "virtio-"
)

var _ csi.NodeServer = &VultrNodeServer{}

type VultrNodeServer struct {
	Driver *VultrDriver
}

func NewVultrNodeDriver(driver *VultrDriver) *VultrNodeServer {
	return &VultrNodeServer{Driver: driver}
}

func (n *VultrNodeServer) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "NodeStageVolume Volume ID must be provided")
	}

	if req.StagingTargetPath == "" {
		return nil, status.Error(codes.InvalidArgument, "NodeStageVolume Staging Target Path must be provided")
	}

	if req.VolumeCapability == nil {
		return nil, status.Error(codes.InvalidArgument, "NodeStageVolume Volume Capability must be provided")
	}

	n.Driver.log.WithFields(logrus.Fields{
		"volume":   req.VolumeId,
		"target":   req.StagingTargetPath,
		"capacity": req.VolumeCapability,
		"method":   "node-stage-method",
	}).Info("node stage volume")

	volumeID, ok := req.GetPublishContext()[n.Driver.publishVolumeID]
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "Could not find the volume id")
	}

	source := getDeviceByPath(volumeID)
	target := req.StagingTargetPath
	mount := req.VolumeCapability.GetMount()
	options := mount.MountFlags

	fsTpe := "ext4"
	if mount.FsType != "" {
		fsTpe = mount.FsType
	}

	//todo format
	// add ability to not reformat disks
	// check if formatted
	// n.Driver.mounter.IsFormatted()
	if err := n.Driver.mounter.Format(source, fsTpe); err != nil {
		n.Driver.log.WithFields(logrus.Fields{
			"source": source,
			"fs":     fsTpe,
			"method": "node-stage-method",
		}).Warn("node stage volume format")
		return nil, status.Error(codes.Internal, err.Error())
	}

	//todo then mount
	// check if mounted
	// n.Driver.mounter.IsMounted()
	if err := n.Driver.mounter.Mount(source, target, fsTpe, options...); err != nil {
		n.Driver.log.WithFields(logrus.Fields{
			"source":  source,
			"target":  target,
			"fs":      fsTpe,
			"options": options,
			"method":  "node-stage-method",
		}).Warn("node stage volume mount")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &csi.NodeStageVolumeResponse{}, nil
}

func (n *VultrNodeServer) NodeUnstageVolume(context.Context, *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	panic("implement me")
}

func (n *VultrNodeServer) NodePublishVolume(context.Context, *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	panic("implement me")
}

func (n *VultrNodeServer) NodeUnpublishVolume(context.Context, *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	panic("implement me")
}

func (n *VultrNodeServer) NodeGetVolumeStats(context.Context, *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (n *VultrNodeServer) NodeExpandVolume(context.Context, *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (n *VultrNodeServer) NodeGetCapabilities(context.Context, *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	nodeCapabilities := []*csi.NodeServiceCapability{
		&csi.NodeServiceCapability{
			Type: &csi.NodeServiceCapability_Rpc{
				Rpc: &csi.NodeServiceCapability_RPC{
					Type: csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME,
				},
			},
		},
	}

	n.Driver.log.WithFields(logrus.Fields{
		"method":       "node-get-capabilities",
		"capabilities": nodeCapabilities,
	})

	return &csi.NodeGetCapabilitiesResponse{
		Capabilities: nodeCapabilities,
	}, nil
}

func (n *VultrNodeServer) NodeGetInfo(ctx context.Context, req *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	n.Driver.log.WithFields(logrus.Fields{
		"method": "node-get-info",
	})

	return &csi.NodeGetInfoResponse{
		NodeId: n.Driver.nodeID,
		AccessibleTopology: &csi.Topology{
			Segments: map[string]string{
				"region": n.Driver.region,
			},
		},
	}, nil
}

func getDeviceByPath(volumeID string) string {
	return filepath.Join(diskPath, fmt.Sprintf("%s%s", diskPrefix, volumeID))
}
