package driver

import (
	"context"
	csi "github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ csi.NodeServer = &VultrNodeServer{}

type VultrNodeServer struct {
	Driver *VultrDriver
}

func NewVultrNodeDriver(driver *VultrDriver) *VultrNodeServer {
	return &VultrNodeServer{Driver: driver}
}

func (n *VultrNodeServer) NodeStageVolume(context.Context, *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	panic("implement me")
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
		NodeId: n.Driver.host,
		AccessibleTopology: &csi.Topology{
			Segments: map[string]string{
				"region": n.Driver.region,
			},
		},
	}, nil
}
