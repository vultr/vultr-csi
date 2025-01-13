// Package vultrstorage is primarily focused on translating CSI storage
// requests into govultr requests for more than one type of storage
// configuration. Currently either block storage or virtual file system storage
package vultrstorage

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/vultr/govultr/v3"
)

const (
	gibiByte int64 = 1073741824

	// Block NVME defaults
	blockNVMEDefaultSize int64 = 10 * gibiByte

	// Block HDD defaults
	blockHDDDefaultSize int64 = 40 * gibiByte
	blockAccessMode           = csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER

	// VFS defaults
	vfsNVMEDefaultSize int64 = 10 * gibiByte
	vfsAccessMode            = csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER
)

// StorageTypes are the available storage types supported by the CSI. Currently
// only "block" & "vfs" are supported
var StorageTypes = []string{"block", "vfs"}

// VultrStorage represents all the relevant data used by the CSI for Vultr
// storages of various types.
type VultrStorage struct {
	ID                string
	Label             string
	Region            string
	BlockType         string
	DiskType          string
	Status            string
	StorageType       string
	SizeGB            int
	AttachedInstances []VultrStorageAttachment
}

// VultrStorageAttachment represents the mount information for a Vultr storage
// attachment.
type VultrStorageAttachment struct {
	NodeID    string
	MountName string
}

// VultrStorageReq represents the general request data used for creating a Vultr storage.
type VultrStorageReq struct {
	Region    string
	SizeGB    int
	Label     string
	BlockType string
	DiskType  string
	Tags      []string // vfs only
}

// VultrStorageUpdateReq represents the general request data used for updating a Vultr storage.
type VultrStorageUpdateReq struct {
	SizeGB int
	Label  string
}

// VultrStorageHandler handles the operations for a VultrStorage of various
// types through its Operations interface.
type VultrStorageHandler struct {
	StorageType  string
	DiskType     string
	DefaultSize  int64
	client       *govultr.Client
	Capabilities []*csi.VolumeCapability
	Operations   interface {
		List(ctx context.Context, options *govultr.ListOptions) ([]VultrStorage, *govultr.Meta, error)
		Get(ctx context.Context, storageID string) (*VultrStorage, error)
		Create(ctx context.Context, req VultrStorageReq) (*VultrStorage, error)
		Update(ctx context.Context, storageID string, req VultrStorageUpdateReq) (*VultrStorage, error)
		Delete(ctx context.Context, storageID string) error
		Attach(ctx context.Context, storageID, instanceID string) error
		Detach(ctx context.Context, storageID, instanceID string) error
	}
}

// NewVultrStorageHandler instantiates a new VultrStorageHandler type and sets
// the Operations interface based on the storageType.
//
// Possible storageTypes are 'block' & 'vfs' for block storage and virtual file
// system storage respectively.
func NewVultrStorageHandler(client *govultr.Client, storageType, diskType string) (*VultrStorageHandler, error) {
	sh := new(VultrStorageHandler)
	sh.client = client
	sh.StorageType = storageType
	sh.DiskType = diskType
	sh.Capabilities = []*csi.VolumeCapability{}

	switch storageType {
	case "block":
		if diskType == "nvme" {
			sh.DefaultSize = blockNVMEDefaultSize
		} else if diskType == "hdd" {
			sh.DefaultSize = blockHDDDefaultSize
		} else {
			sh.DefaultSize = 0
		}

		sh.Operations = &VultrBlockStorageHandler{client}
		sh.Capabilities = append(sh.Capabilities, &csi.VolumeCapability{
			AccessMode: &csi.VolumeCapability_AccessMode{
				Mode: blockAccessMode,
			},
		})
		return sh, nil

	case "vfs":
		if diskType == "nvme" {
			sh.DefaultSize = vfsNVMEDefaultSize
		} else {
			sh.DefaultSize = 0
		}

		sh.Operations = &VultrVFSStorageHandler{client}
		sh.Capabilities = append(sh.Capabilities, &csi.VolumeCapability{
			AccessMode: &csi.VolumeCapability_AccessMode{
				Mode: vfsAccessMode,
			},
		})
		return sh, nil
	}

	return nil, fmt.Errorf("unable to instantiate a new storage handler : invalid storage type")
}

// FindVultrStorageHandlerByID performs a lookup of available storage types and
// returns the appropriate handler to use with the storage
func FindVultrStorageHandlerByID(ctx context.Context, client *govultr.Client, storageID string) (*VultrStorageHandler, error) {
	if storageID == "" {
		return nil, fmt.Errorf("missing storage ID")
	}

	for _, storageType := range StorageTypes {
		sh, err := NewVultrStorageHandler(client, storageType, "")
		if err != nil {
			return nil, fmt.Errorf("FindVultrStorageHandlerByID cannot initialize vultr storage handler. %v", err)
		}

		storage, err := sh.Operations.Get(ctx, storageID)
		if err != nil {
			if strings.Contains(err.Error(), "Invalid block storage ID") ||
				strings.Contains(err.Error(), "Subscription ID Not Found") {
				continue
			}

			// some other error
			return nil, fmt.Errorf("FindVultrStorageHandlerByID could not retrieve storage: %v", err)
		}

		if storage != nil {
			return sh, nil
		}
	}

	return nil, fmt.Errorf("storage not found : %v", storageID)
}

// ListAllStorages retrieves the list results of available storage types and
// returns the VultrStorage objects for further processing.
func ListAllStorages(ctx context.Context, client *govultr.Client) ([]VultrStorage, error) {
	var allStorages []VultrStorage
	for _, storageType := range StorageTypes {
		sh, err := NewVultrStorageHandler(client, storageType, "")
		if err != nil {
			return nil, fmt.Errorf("ListAllStorages cannot initialize vultr storage handler. %v", err)
		}

		listOptions := &govultr.ListOptions{}

		for {
			storages, meta, err := sh.Operations.List(ctx, listOptions)
			if err != nil {
				return nil, fmt.Errorf("ListVolumes cannot retrieve list of volumes. %v", err)
			}

			allStorages = append(allStorages, storages...)

			if meta.Links.Next != "" {
				listOptions.Cursor = meta.Links.Next
				continue
			}
			break
		}
	}

	return allStorages, nil
}

// Block Storage ==============================================================

// VultrBlockStorageHandler implements the Operations interface on the
// VultrStorageHandler and performs those operations for block storages.
type VultrBlockStorageHandler struct {
	client *govultr.Client
}

// List wraps the govultr List and converts the response for block storage.
func (v *VultrBlockStorageHandler) List(ctx context.Context, options *govultr.ListOptions) ([]VultrStorage, *govultr.Meta, error) {
	bss, meta, _, err := v.client.BlockStorage.List(ctx, options)
	if err != nil {
		return nil, nil, err
	}

	var vss []VultrStorage
	for i := range bss {
		vs, err := convertFromBlock(&bss[i])
		if err != nil {
			continue
		}

		vss = append(vss, *vs)
	}

	return vss, meta, nil
}

// Get wraps the govultr Get function and converts the responses for block
// storage.
func (v *VultrBlockStorageHandler) Get(ctx context.Context, blockID string) (*VultrStorage, error) {
	bs, _, err := v.client.BlockStorage.Get(ctx, blockID)
	if err != nil {
		return nil, fmt.Errorf("storage handler unable to retrieve block storage : %s", err)
	}

	vs, err := convertFromBlock(bs)
	if err != nil {
		return nil, fmt.Errorf("storage handler unable to convert block storage data on get : %s", err)
	}

	return vs, nil
}

// Create wraps the govultr Create function and converts the response for block
// storage.
func (v *VultrBlockStorageHandler) Create(ctx context.Context, req VultrStorageReq) (*VultrStorage, error) {
	bsReq := new(govultr.BlockStorageCreate)
	bsReq.Region = req.Region
	bsReq.Label = req.Label
	bsReq.SizeGB = req.SizeGB

	switch req.DiskType {
	case "hdd":
		bsReq.BlockType = "storage_opt"
	case "nvme":
		bsReq.BlockType = "high_perf"
	}

	bs, _, err := v.client.BlockStorage.Create(ctx, bsReq)
	if err != nil {
		return nil, fmt.Errorf("storage handler unable to create block storage : %s", err)
	}

	vs, err := convertFromBlock(bs)
	if err != nil {
		return nil, fmt.Errorf("storage handler unable to convert block storage data on create : %s", err)
	}

	return vs, nil
}

// Update wraps the govultr Update function and converts the response for block
// storage.
func (v *VultrBlockStorageHandler) Update(ctx context.Context, storageID string, req VultrStorageUpdateReq) (*VultrStorage, error) {
	bsReq := new(govultr.BlockStorageUpdate)
	bsReq.Label = req.Label
	bsReq.SizeGB = req.SizeGB

	if err := v.client.BlockStorage.Update(ctx, storageID, bsReq); err != nil {
		return nil, fmt.Errorf("storage handler unable to update block storage : %s", err)
	}

	// block storage doesn't return anything on update so let's get that for
	// consistency
	vs, err := v.Get(ctx, storageID)
	if err != nil {
		return nil, fmt.Errorf("storage handler unable to retrieve block storage data on update : %s", err)
	}

	return vs, nil
}

// Delete wraps the govultr Delete function for block storage.
func (v *VultrBlockStorageHandler) Delete(ctx context.Context, storageID string) error {
	if err := v.client.BlockStorage.Delete(ctx, storageID); err != nil {
		return fmt.Errorf("storage handler unable to delete block storage : %s", err)
	}

	return nil
}

// Attach wraps the govultr Attach function for block storage.
func (v *VultrBlockStorageHandler) Attach(ctx context.Context, storageID, instanceID string) error {
	attReq := govultr.BlockStorageAttach{
		InstanceID: instanceID,
		Live:       govultr.BoolToBoolPtr(true),
	}

	if err := v.client.BlockStorage.Attach(ctx, storageID, &attReq); err != nil {
		return fmt.Errorf("storage handler unable to attach block storage : %s", err)
	}

	return nil
}

// Detach wraps the govultr Detach function for block storage.
func (v *VultrBlockStorageHandler) Detach(ctx context.Context, storageID, instanceID string) error {
	detReq := govultr.BlockStorageDetach{
		Live: govultr.BoolToBoolPtr(true),
	}

	if err := v.client.BlockStorage.Detach(ctx, storageID, &detReq); err != nil {
		return fmt.Errorf("storage handler unable to detach block storage : %s", err)
	}

	return nil
}

// VFS Storage ================================================================

// VultrVFSStorageHandler implements the Operations interface on the
// VultrStorageHandler and performs those operations for VFS storages.
type VultrVFSStorageHandler struct {
	client *govultr.Client
}

// List wraps the govultr List and converts the response for VFS storage.
func (v *VultrVFSStorageHandler) List(ctx context.Context, options *govultr.ListOptions) ([]VultrStorage, *govultr.Meta, error) {
	vfss, meta, _, err := v.client.VirtualFileSystemStorage.List(ctx, options)
	if err != nil {
		return nil, nil, fmt.Errorf("storage handler unable to retrieve vfs storage list : %s", err)
	}

	// List checks in CSI do not check for attached instances so skip the lookup

	var s []VultrStorage
	for i := range vfss {
		vfs, err := convertFromVFS(&vfss[i], nil)
		if err != nil {
			continue
		}

		s = append(s, *vfs)
	}

	return s, meta, nil
}

// Get wraps the govultr Get and AttachmentList functions and converts the
// responses for VFS storage.
func (v *VultrVFSStorageHandler) Get(ctx context.Context, storageID string) (*VultrStorage, error) {
	vfs, _, err := v.client.VirtualFileSystemStorage.Get(ctx, storageID)
	if err != nil {
		return nil, fmt.Errorf("storage handler unable to retrieve vfs storage : %s", err)
	}

	// VFS does not include attached instance in the get, must lookup separately
	attached, _, err := v.client.VirtualFileSystemStorage.AttachmentList(ctx, storageID)
	if err != nil {
		return nil, fmt.Errorf("storage handler unable to lookup attached instances for vfs storage : %s", err)
	}

	vs, err := convertFromVFS(vfs, attached)
	if err != nil {
		return nil, fmt.Errorf("storage handler unable to convert vfs storage data on get : %s", err)
	}

	return vs, nil
}

// Create wraps the govultr Create function and converts the respons for a VFS
// storage.
func (v *VultrVFSStorageHandler) Create(ctx context.Context, req VultrStorageReq) (*VultrStorage, error) {
	vfsReq := new(govultr.VirtualFileSystemStorageReq)
	vfsReq.Region = req.Region
	vfsReq.Label = req.Label
	vfsReq.StorageSize.SizeGB = req.SizeGB
	vfsReq.DiskType = req.DiskType
	vfsReq.Tags = req.Tags

	vfs, _, err := v.client.VirtualFileSystemStorage.Create(ctx, vfsReq)
	if err != nil {
		return nil, fmt.Errorf("storage handler unable to create vfs storage : %s", err)
	}

	// nothing can be attached at creation so pass nil
	vs, err := convertFromVFS(vfs, nil)
	if err != nil {
		return nil, fmt.Errorf("storage handler unable to convert vfs storage data on create : %s", err)
	}

	return vs, nil
}

// Update wraps the govultr Update function and converts the response for VFS
// storage.
func (v *VultrVFSStorageHandler) Update(ctx context.Context, storageID string, req VultrStorageUpdateReq) (*VultrStorage, error) {
	vfsReq := new(govultr.VirtualFileSystemStorageUpdateReq)
	vfsReq.Label = req.Label
	vfsReq.StorageSize.SizeGB = req.SizeGB

	vfs, _, err := v.client.VirtualFileSystemStorage.Update(ctx, storageID, vfsReq)
	if err != nil {
		return nil, fmt.Errorf("storage handler unable to update vfs storage : %s", err)
	}

	// VFS does not include attached instance in the update, must lookup separately
	attached, _, err := v.client.VirtualFileSystemStorage.AttachmentList(ctx, storageID)
	if err != nil {
		return nil, fmt.Errorf("storage handler unable to lookup attached instances for vfs storage update : %s", err)
	}

	vs, err := convertFromVFS(vfs, attached)
	if err != nil {
		return nil, fmt.Errorf("storage handler unable to convert vfs storage data on update : %s", err)
	}

	return vs, nil
}

// Delete wraps the govultr Delete function for vfs storage.
func (v *VultrVFSStorageHandler) Delete(ctx context.Context, storageID string) error {
	if err := v.client.VirtualFileSystemStorage.Delete(ctx, storageID); err != nil {
		return fmt.Errorf("storage handler unable to delete vfs storage : %s", err)
	}

	return nil
}

// Attach wraps the govultr Attach function for VFS storage.
func (v *VultrVFSStorageHandler) Attach(ctx context.Context, storageID, instanceID string) error {
	if _, _, err := v.client.VirtualFileSystemStorage.Attach(ctx, storageID, instanceID); err != nil {
		return fmt.Errorf("storage handler unable to attach vfs storage : %s", err)
	}

	return nil
}

// Detach wraps the govultr Detach function for vfs storage.
func (v *VultrVFSStorageHandler) Detach(ctx context.Context, storageID, instanceID string) error {
	if err := v.client.VirtualFileSystemStorage.Detach(ctx, storageID, instanceID); err != nil {
		return fmt.Errorf("storage handler unable to detach vfs storage : %s", err)
	}

	return nil
}

// Data mapping functions ==========================================================

func convertFromVFS(vfs *govultr.VirtualFileSystemStorage, attached []govultr.VirtualFileSystemStorageAttachment) (*VultrStorage, error) {
	if vfs == nil {
		return nil, fmt.Errorf("vfs storage is empty")
	}

	vs := new(VultrStorage)
	vs.ID = vfs.ID
	vs.Label = vfs.Label
	vs.SizeGB = vfs.StorageSize.SizeGB
	vs.Region = vfs.Region
	vs.DiskType = vfs.DiskType
	vs.Status = vfs.Status
	vs.StorageType = "vfs"

	// Not relevant to vfs
	vs.BlockType = ""

	for i := range attached {
		vs.AttachedInstances = append(vs.AttachedInstances, VultrStorageAttachment{
			NodeID:    attached[i].TargetID,
			MountName: strconv.Itoa(attached[i].MountTag),
		})
	}

	return vs, nil
}

func convertFromBlock(bs *govultr.BlockStorage) (*VultrStorage, error) {
	if bs == nil {
		return nil, fmt.Errorf("block storage is empty")
	}

	vs := new(VultrStorage)
	vs.ID = bs.ID
	vs.Label = bs.Label
	vs.SizeGB = bs.SizeGB
	vs.Region = bs.Region
	vs.BlockType = bs.BlockType
	vs.StorageType = "block"

	if bs.BlockType == "high_perf" {
		vs.DiskType = "nvme"
	} else if bs.BlockType == "storage_opt" {
		vs.DiskType = "hdd"
	} else {
		vs.DiskType = ""
	}
	vs.Status = bs.Status

	// Block storage can only be attached to one instance
	if bs.AttachedToInstance != "" {
		vs.AttachedInstances = append(vs.AttachedInstances, VultrStorageAttachment{
			NodeID:    bs.AttachedToInstance,
			MountName: bs.MountID,
		})
	}

	return vs, nil
}
