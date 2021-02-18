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
	}).Info("Node Stage Volume: called")

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

	formatted, err := n.Driver.mounter.IsFormatted(source)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot verify if formatted: %v", err.Error())
	}

	if !formatted {
		if err = n.Driver.mounter.Format(source, fsTpe); err != nil {
			n.Driver.log.WithFields(logrus.Fields{
				"source": source,
				"fs":     fsTpe,
				"method": "node-stage-method",
			}).Warn("node stage volume format")
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	mounted, err := n.Driver.mounter.IsMounted(target)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if !mounted {
		if err := n.Driver.mounter.Mount(source, target, fsTpe, options...); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	n.Driver.log.Info("Node Stage Volume: volume staged")
	return &csi.NodeStageVolumeResponse{}, nil
}

func (n *VultrNodeServer) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
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

	mounted, err := n.Driver.mounter.IsMounted(req.StagingTargetPath)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot verify mount status for %v, %v", req.StagingTargetPath, err.Error())
	}

	if mounted {
		err := n.Driver.mounter.UnMount(req.StagingTargetPath)
		if err != nil {
			return nil, err
		}
	}

	n.Driver.log.Info("Node Unstage Volume: volume unstaged")
	return &csi.NodeUnstageVolumeResponse{}, nil
}

func (n *VultrNodeServer) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
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
	for _, flag := range mnt.MountFlags {
		options = append(options, flag)
	}

	fsType := "ext4"
	if mnt.FsType != "" {
		fsType = mnt.FsType
	}

	mounted, err := n.Driver.mounter.IsMounted(req.TargetPath)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot verify mount status for %v, %v", req.StagingTargetPath, err.Error())
	}

	if !mounted {
		err := n.Driver.mounter.Mount(req.StagingTargetPath, req.TargetPath, fsType, options...)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	n.Driver.log.Info("Node Publish Volume: published")
	return &csi.NodePublishVolumeResponse{}, nil
}

func (n *VultrNodeServer) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
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

	mounted, err := n.Driver.mounter.IsMounted(req.TargetPath)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot verify mount status for %v, %v", req.TargetPath, err.Error())
	}

	if mounted {
		err := n.Driver.mounter.UnMount(req.TargetPath)
		if err != nil {
			return nil, err
		}
	}

	n.Driver.log.Info("Node Publish Volume: unpublished")
	return &csi.NodeUnpublishVolumeResponse{}, nil
}

func (n *VultrNodeServer) NodeGetVolumeStats(context.Context, *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (n *VultrNodeServer) NodeExpandVolume(context.Context, *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (n *VultrNodeServer) NodeGetCapabilities(context.Context, *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	nodeCapabilities := []*csi.NodeServiceCapability{
		{
			Type: &csi.NodeServiceCapability_Rpc{
				Rpc: &csi.NodeServiceCapability_RPC{
					Type: csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME,
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

func (n *VultrNodeServer) NodeGetInfo(ctx context.Context, req *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	n.Driver.log.WithFields(logrus.Fields{}).Info("Node Get Info: called")

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
