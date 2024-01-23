package driver

import (
	"context"
	"fmt"
	"path/filepath"

	csi "github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	diskPath   = "/dev/disk/by-id"
	diskPrefix = "virtio-"

	maxVolumesPerNode = 11

	volumeModeBlock      = "block"
	volumeModeFilesystem = "filesystem"
)

var _ csi.NodeServer = &VultrNodeServer{}

// VultrNodeServer type provides the VultrDriver
type VultrNodeServer struct {
	Driver *VultrDriver
}

// NewVultrNodeDriver provides a VultrNodeServer
func NewVultrNodeDriver(driver *VultrDriver) *VultrNodeServer {
	return &VultrNodeServer{Driver: driver}
}

// NodeStageVolume provides stages the node volume
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
	}).Info("Node Stage Volume: called")

	volumeID, ok := req.GetPublishContext()[n.Driver.mountID]
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

	if err := n.Driver.mounter.FormatAndMount(source, target, fsTpe, options); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	n.Driver.log.Info("Node Stage Volume: volume staged")
	return &csi.NodeStageVolumeResponse{}, nil
}

// NodeUnstageVolume provides the node volume unstage functionality
func (n *VultrNodeServer) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) { //nolint:dupl,lll
	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "VolumeID must be provided")
	}

	if req.StagingTargetPath == "" {
		return nil, status.Error(codes.InvalidArgument, "Staging Target Path must be provided")
	}

	n.Driver.log.WithFields(logrus.Fields{
		"volume-id":           req.VolumeId,
		"staging-target-path": req.StagingTargetPath,
	}).Info("Node Unstage Volume: called")

	err := n.Driver.mounter.Unmount(req.StagingTargetPath)
	if err != nil {
		return nil, err
	}

	n.Driver.log.Info("Node Unstage Volume: volume unstaged")
	return &csi.NodeUnstageVolumeResponse{}, nil
}

// NodePublishVolume allows the volume publish
func (n *VultrNodeServer) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) { //nolint:lll
	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "VolumeID must be provided")
	}

	if req.StagingTargetPath == "" {
		return nil, status.Error(codes.InvalidArgument, "Staging Target Path must be provided")
	}

	if req.TargetPath == "" {
		return nil, status.Error(codes.InvalidArgument, "Target Path must be provided")
	}

	log := n.Driver.log.WithFields(logrus.Fields{
		"volume_id":           req.VolumeId,
		"staging_target_path": req.StagingTargetPath,
		"target_path":         req.TargetPath,
	})
	log.Info("Node Publish Volume: called")

	options := []string{"bind"}
	if req.Readonly {
		options = append(options, "ro")
	}

	mnt := req.VolumeCapability.GetMount()
	options = append(options, mnt.MountFlags...)

	fsType := "ext4"
	if mnt.FsType != "" {
		fsType = mnt.FsType
	}

	err := n.Driver.mounter.Mount(req.StagingTargetPath, req.TargetPath, fsType, options)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	n.Driver.log.Info("Node Publish Volume: published")
	return &csi.NodePublishVolumeResponse{}, nil
}

// NodeUnpublishVolume allows the volume to be unpublished
func (n *VultrNodeServer) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) { //nolint:dupl,lll
	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "VolumeID must be provided")
	}

	if req.TargetPath == "" {
		return nil, status.Error(codes.InvalidArgument, "Target Path must be provided")
	}

	n.Driver.log.WithFields(logrus.Fields{
		"volume-id":   req.VolumeId,
		"target-path": req.TargetPath,
	}).Info("Node Unpublish Volume: called")

	err := n.Driver.mounter.Unmount(req.TargetPath)
	if err != nil {
		return nil, err
	}

	n.Driver.log.Info("Node Publish Volume: unpublished")
	return &csi.NodeUnpublishVolumeResponse{}, nil
}

// NodeGetVolumeStats provides the volume stats
func (n *VultrNodeServer) NodeGetVolumeStats(ctx context.Context, req *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) { //nolint:lll
	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "NodeGetVolumeStats Volume ID must be provided")
	}

	volumePath := req.VolumePath
	if volumePath == "" {
		return nil, status.Error(codes.InvalidArgument, "NodeGetVolumeStats Volume Path must be provided")
	}

	log := n.Driver.log.WithFields(logrus.Fields{
		"volume_id":   req.VolumeId,
		"volume_path": req.VolumePath,
		"method":      "node_get_volume_stats",
	})
	log.Info("node get volume stats called")

	mounted, err := n.Driver.vMounter.IsMounted(volumePath)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check if volume path %q is mounted: %s", volumePath, err)
	}

	if !mounted {
		return nil, status.Errorf(codes.NotFound, "volume path %q is not mounted", volumePath)
	}

	isBlock, err := n.Driver.vMounter.IsBlockDevice(volumePath)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to determine if %q is block device: %s", volumePath, err)
	}

	stats, err := n.Driver.vMounter.GetStatistics(volumePath)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve capacity statistics for volume path %q: %s", volumePath, err)
	}

	// only can retrieve total capacity for a block device
	if isBlock {
		log.WithFields(logrus.Fields{
			"volume_mode": volumeModeBlock,
			"bytes_total": stats.totalBytes,
		}).Info("node capacity statistics retrieved")

		return &csi.NodeGetVolumeStatsResponse{
			Usage: []*csi.VolumeUsage{
				{
					Unit:  csi.VolumeUsage_BYTES,
					Total: stats.totalBytes,
				},
			},
		}, nil
	}

	log.WithFields(logrus.Fields{
		"volume_mode":      volumeModeFilesystem,
		"bytes_available":  stats.availableBytes,
		"bytes_total":      stats.totalBytes,
		"bytes_used":       stats.usedBytes,
		"inodes_available": stats.availableInodes,
		"inodes_total":     stats.totalInodes,
		"inodes_used":      stats.usedInodes,
	}).Info("node capacity statistics retrieved")

	return &csi.NodeGetVolumeStatsResponse{
		Usage: []*csi.VolumeUsage{
			{
				Available: stats.availableBytes,
				Total:     stats.totalBytes,
				Used:      stats.usedBytes,
				Unit:      csi.VolumeUsage_BYTES,
			},
			{
				Available: stats.availableInodes,
				Total:     stats.totalInodes,
				Used:      stats.usedInodes,
				Unit:      csi.VolumeUsage_INODES,
			},
		},
	}, nil
}

// NodeExpandVolume provides the node volume expansion
func (n *VultrNodeServer) NodeExpandVolume(ctx context.Context, req *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	n.Driver.log.WithFields(logrus.Fields{
		"required_bytes": req.CapacityRange.RequiredBytes,
	}).Info("Node Expand Volume: called")

	return &csi.NodeExpandVolumeResponse{
		CapacityBytes: req.CapacityRange.RequiredBytes,
	}, nil
}

// NodeGetCapabilities provides the node capabilities
func (n *VultrNodeServer) NodeGetCapabilities(context.Context, *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	nodeCapabilities := []*csi.NodeServiceCapability{
		{
			Type: &csi.NodeServiceCapability_Rpc{
				Rpc: &csi.NodeServiceCapability_RPC{
					Type: csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME,
				},
			},
		},
		{
			Type: &csi.NodeServiceCapability_Rpc{
				Rpc: &csi.NodeServiceCapability_RPC{
					Type: csi.NodeServiceCapability_RPC_GET_VOLUME_STATS,
				},
			},
		},
		{
			Type: &csi.NodeServiceCapability_Rpc{
				Rpc: &csi.NodeServiceCapability_RPC{
					Type: csi.NodeServiceCapability_RPC_EXPAND_VOLUME,
				},
			},
		},
	}

	n.Driver.log.WithFields(logrus.Fields{
		"capabilities": nodeCapabilities,
	}).Info("Node Get Capabilities: called")

	return &csi.NodeGetCapabilitiesResponse{
		Capabilities: nodeCapabilities,
	}, nil
}

// NodeGetInfo provides the node info
func (n *VultrNodeServer) NodeGetInfo(ctx context.Context, req *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	n.Driver.log.WithFields(logrus.Fields{}).Info("Node Get Info: called")

	return &csi.NodeGetInfoResponse{
		NodeId:            n.Driver.nodeID,
		MaxVolumesPerNode: maxVolumesPerNode,
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
