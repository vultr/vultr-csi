package driver

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
	"k8s.io/mount-utils"

	"github.com/container-storage-interface/spec/lib/go/csi"
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
	mountBlk := req.VolumeCapability.GetMount()
	options := mountBlk.MountFlags

	fsTpe := "ext4"
	if mountBlk.FsType != "" {
		fsTpe = mountBlk.FsType
	}

	n.Driver.log.WithFields(logrus.Fields{
		"volume":   req.VolumeId,
		"target":   req.StagingTargetPath,
		"capacity": req.VolumeCapability,
	}).Infof("Node Stage Volume: creating directory target %s\n", target)
	err := os.MkdirAll(target, mkDirMode)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	n.Driver.log.WithFields(logrus.Fields{
		"volume":   req.VolumeId,
		"target":   req.StagingTargetPath,
		"capacity": req.VolumeCapability,
	}).Infof("Node Stage Volume: directory created for target %s\n", target)

	n.Driver.log.WithFields(logrus.Fields{
		"volume":   req.VolumeId,
		"target":   req.StagingTargetPath,
		"capacity": req.VolumeCapability,
	}).Info("Node Stage Volume: attempting format and mount")

	if err := n.Driver.mounter.FormatAndMount(source, target, fsTpe, options); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if _, err := os.Stat(source); err == nil {
		needResize, err := n.Driver.resizer.NeedResize(source, target)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "could not determine if volume %q needs to be resized: %v", req.VolumeId, err)
		}

		if needResize {
			n.Driver.log.WithFields(logrus.Fields{
				"volume":   req.VolumeId,
				"target":   req.StagingTargetPath,
				"capacity": req.VolumeCapability,
			}).Info("Node Stage Volume: resizing volume")

			if _, err := n.Driver.resizer.Resize(source, target); err != nil {
				return nil, status.Errorf(codes.Internal, "could not resize volume %q:  %v", req.VolumeId, err)
			}
		}
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

	err := os.MkdirAll(req.TargetPath, mkDirMode)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = n.Driver.mounter.Mount(req.StagingTargetPath, req.TargetPath, fsType, options)
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

	statfs := &unix.Statfs_t{}
	err := unix.Statfs(volumePath, statfs)
	if err != nil {
		return nil, err
	}

	availableBytes := int64(statfs.Bavail) * int64(statfs.Bsize)                    //nolint
	usedBytes := (int64(statfs.Blocks) - int64(statfs.Bfree)) * int64(statfs.Bsize) //nolint
	totalBytes := int64(statfs.Blocks) * int64(statfs.Bsize)                        //nolint
	totalInodes := int64(statfs.Files)
	availableInodes := int64(statfs.Ffree)
	usedInodes := totalInodes - availableInodes

	log.WithFields(logrus.Fields{
		"volume_mode":      volumeModeFilesystem,
		"bytes_available":  availableBytes,
		"bytes_total":      totalBytes,
		"bytes_used":       usedBytes,
		"inodes_available": availableInodes,
		"inodes_total":     totalInodes,
		"inodes_used":      usedInodes,
	}).Info("node capacity statistics retrieved")

	return &csi.NodeGetVolumeStatsResponse{
		Usage: []*csi.VolumeUsage{
			{
				Available: availableBytes,
				Total:     totalBytes,
				Used:      usedBytes,
				Unit:      csi.VolumeUsage_BYTES,
			},
			{
				Available: availableInodes,
				Total:     totalInodes,
				Used:      usedInodes,
				Unit:      csi.VolumeUsage_INODES,
			},
		},
	}, nil
}

// NodeExpandVolume provides the node volume expansion
func (n *VultrNodeServer) NodeExpandVolume(ctx context.Context, req *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	log := n.Driver.log.WithFields(logrus.Fields{
		"volume_id":   req.VolumeId,
		"volume_path": req.VolumePath,
		"method":      "NodeExpandVolume",
	})

	n.Driver.log.WithFields(logrus.Fields{
		"required_bytes": req.CapacityRange.RequiredBytes,
	}).Info("Node Expand Volume: called")

	devicePath, _, err := mount.GetDeviceNameFromMount(mount.New(""), req.VolumePath)
	if err != nil {
		log.Infof("failed to determine mount path for %s: %s", req.VolumePath, err)
		return nil, fmt.Errorf("failed to determine mount path for %s: %s", req.VolumePath, err)
	}

	log.Infof("attempting to resize devicepath: %s", devicePath)

	if _, err := n.Driver.resizer.Resize(devicePath, req.VolumePath); err != nil {
		log.Infof("failed to resize volume: %s", err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to resize volume: %s", err))
	}

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
