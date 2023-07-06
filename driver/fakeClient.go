package driver

import (
	"context"

	"github.com/vultr/govultr/v2"
)

func newFakeClient() *govultr.Client {
	fakeInstance := FakeInstance{client: nil}
	fakeBlockStorage := fakeBS{client: nil}

	return &govultr.Client{
		Instance:     &fakeInstance,
		BlockStorage: &fakeBlockStorage,
	}
}

func newFakeBS() *govultr.BlockStorage {
	return &govultr.BlockStorage{
		ID:                 "c56c7b6e-15c2-445e-9a5d-1063ab5828ec",
		DateCreated:        "",
		Cost:               10,
		Status:             "active",
		SizeGB:             10,
		Region:             "ewr",
		AttachedToInstance: "245bb2fe-b55c-44a0-9a1e-ab80e4b5f088",
		Label:              "test-bs",
		MountID:            "c56c7b6e-15c2-445e-9a5d-1063ab5828ec",
	}
}

type fakeBS struct {
	client *govultr.Client
}

func (f *fakeBS) Create(ctx context.Context, blockReq *govultr.BlockStorageCreate) (*govultr.BlockStorage, error) {
	return newFakeBS(), nil
}

func (f *fakeBS) Get(ctx context.Context, blockID string) (*govultr.BlockStorage, error) {
	return newFakeBS(), nil
}

func (f *fakeBS) Update(ctx context.Context, blockID string, blockReq *govultr.BlockStorageUpdate) error {
	panic("implement me")
}

func (f *fakeBS) Delete(ctx context.Context, blockID string) error {
	return nil
}

func (f *fakeBS) List(ctx context.Context, options *govultr.ListOptions) ([]govultr.BlockStorage, *govultr.Meta, error) {
	return []govultr.BlockStorage{
			{
				ID:                 "c56c7b6e-15c2-445e-9a5d-1063ab5828ec",
				DateCreated:        "",
				Cost:               10,
				Status:             "active",
				SizeGB:             10,
				Region:             "ewr",
				AttachedToInstance: "245bb2fe-b55c-44a0-9a1e-ab80e4b5f088",
				Label:              "test-bs",
				MountID:            "c56c7b6e-15c2-445e-9a5d-1063ab5828ec",
			},
			{
				ID:                 "bda4f333-bfd7-477b-84c2-e4df0ec9e5bf",
				DateCreated:        "",
				Cost:               20,
				Status:             "active",
				SizeGB:             20,
				Region:             "ewr",
				AttachedToInstance: "b9d23eb3-1880-4746-acc7-f1ef56565320",
				Label:              "test-bs2",
				MountID:            "b9d23eb3-1880-4746-acc7-f1ef56565320",
			},
		}, &govultr.Meta{
			Total: 0,
			Links: &govultr.Links{
				Next: "",
				Prev: "",
			},
		}, nil
}

func (f *fakeBS) Attach(ctx context.Context, blockID string, attach *govultr.BlockStorageAttach) error {
	panic("implement me")
}

func (f *fakeBS) Detach(ctx context.Context, blockID string, detach *govultr.BlockStorageDetach) error {
	list, _, err := f.List(ctx, nil)
	if err != nil {
		return err
	}

	for i := range list {
		if list[i].ID == blockID {
			list[i].AttachedToInstance = ""
		}
	}

	return nil
}

// FakeInstance returns the client
type FakeInstance struct {
	client *govultr.Client
}

// Create is not implemented
func (f *FakeInstance) Create(ctx context.Context, instanceReq *govultr.InstanceCreateReq) (*govultr.Instance, error) {
	panic("implement me")
}

// Get returns an instance struct
func (f *FakeInstance) Get(ctx context.Context, instanceID string) (*govultr.Instance, error) {
	return &govultr.Instance{
		ID:           "94cf529e-796c-44c0-8a18-6e0be753f155",
		MainIP:       "149.28.225.110",
		VCPUCount:    4,
		Region:       "ewr",
		Status:       "running",
		NetmaskV4:    "255.255.254.0",
		GatewayV4:    "149.28.224.1",
		PowerStatus:  "",
		ServerStatus: "",
		Plan:         "vc2-4c-8gb",
		Label:        "csi-test",
		InternalIP:   "10.1.95.4",
	}, nil
}

// Update is not implemented
func (f *FakeInstance) Update(ctx context.Context, instanceID string, instanceReq *govultr.InstanceUpdateReq) (*govultr.Instance, error) {
	panic("implement me")
}

// Delete jis not implemented
func (f *FakeInstance) Delete(ctx context.Context, instanceID string) error {
	panic("implement me")
}

// List is not implemented
func (f *FakeInstance) List(ctx context.Context, options *govultr.ListOptions) ([]govultr.Instance, *govultr.Meta, error) {
	panic("implement me")
}

// Start is not implemented
func (f *FakeInstance) Start(ctx context.Context, instanceID string) error {
	panic("implement me")
}

// Halt is not implemented
func (f *FakeInstance) Halt(ctx context.Context, instanceID string) error {
	panic("implement me")
}

// Reboot is not implemented
func (f *FakeInstance) Reboot(ctx context.Context, instanceID string) error {
	panic("implement me")
}

// Reinstall is not implemented
func (f *FakeInstance) Reinstall(ctx context.Context, instanceID string, reinstallReq *govultr.ReinstallReq) (*govultr.Instance, error) {
	panic("implement me")
}

// MassStart is not implemented
func (f *FakeInstance) MassStart(ctx context.Context, instanceList []string) error {
	panic("implement me")
}

// MassHalt is not implemented
func (f *FakeInstance) MassHalt(ctx context.Context, instanceList []string) error {
	panic("implement me")
}

// MassReboot is not implemented
func (f *FakeInstance) MassReboot(ctx context.Context, instanceList []string) error {
	panic("implement me")
}

// Restore is not implemented
func (f *FakeInstance) Restore(ctx context.Context, instanceID string, restoreReq *govultr.RestoreReq) error {
	panic("implement me")
}

// GetBandwidth is not implemented
func (f *FakeInstance) GetBandwidth(ctx context.Context, instanceID string) (*govultr.Bandwidth, error) {
	panic("implement me")
}

// GetNeighbors is not implemented
func (f *FakeInstance) GetNeighbors(ctx context.Context, instanceID string) (*govultr.Neighbors, error) {
	panic("implement me")
}

// ListPrivateNetworks is not implemented
func (f *FakeInstance) ListPrivateNetworks(ctx context.Context, instanceID string, options *govultr.ListOptions) ([]govultr.PrivateNetwork, *govultr.Meta, error) { //nolint:lll
	panic("implement me")
}

// ListVPCInfo is not implemented
func (f *FakeInstance) ListVPCInfo(ctx context.Context, instanceID string, options *govultr.ListOptions) ([]govultr.VPCInfo, *govultr.Meta, error) { //nolint:lll
	panic("implement me")
}

// AttachPrivateNetwork is not implemented
func (f *FakeInstance) AttachPrivateNetwork(ctx context.Context, instanceID, networkID string) error {
	panic("implement me")
}

// AttachVPC is not implemented
func (f *FakeInstance) AttachVPC(ctx context.Context, instanceID, networkID string) error {
	panic("implement me")
}

// DetachPrivateNetwork is not implemented
func (f *FakeInstance) DetachPrivateNetwork(ctx context.Context, instanceID, networkID string) error {
	panic("implement me")
}

// DetachVPC is not implemented
func (f *FakeInstance) DetachVPC(ctx context.Context, instanceID, networkID string) error {
	panic("implement me")
}

// ISOStatus is not implemented
func (f *FakeInstance) ISOStatus(ctx context.Context, instanceID string) (*govultr.Iso, error) {
	panic("implement me")
}

// AttachISO is not implemented
func (f *FakeInstance) AttachISO(ctx context.Context, instanceID, isoID string) error {
	panic("implement me")
}

// DetachISO is not implemented
func (f *FakeInstance) DetachISO(ctx context.Context, instanceID string) error {
	panic("implement me")
}

// GetBackupSchedule is not implemented
func (f *FakeInstance) GetBackupSchedule(ctx context.Context, instanceID string) (*govultr.BackupSchedule, error) {
	panic("implement me")
}

// SetBackupSchedule is not implemented
func (f *FakeInstance) SetBackupSchedule(ctx context.Context, instanceID string, backup *govultr.BackupScheduleReq) error {
	panic("implement me")
}

// CreateIPv4 is not implemented
func (f *FakeInstance) CreateIPv4(ctx context.Context, instanceID string, reboot *bool) (*govultr.IPv4, error) {
	panic("implement me")
}

// ListIPv4 is not implemented
func (f *FakeInstance) ListIPv4(ctx context.Context, instanceID string, option *govultr.ListOptions) ([]govultr.IPv4, *govultr.Meta, error) { //nolint:lll
	panic("implement me")
}

// DeleteIPv4 is not implemented
func (f *FakeInstance) DeleteIPv4(ctx context.Context, instanceID, ip string) error {
	panic("implement me")
}

// ListIPv6 is not implemented
func (f *FakeInstance) ListIPv6(ctx context.Context, instanceID string, option *govultr.ListOptions) ([]govultr.IPv6, *govultr.Meta, error) { //nolint:lll
	panic("implement me")
}

// CreateReverseIPv6 is not implemented
func (f *FakeInstance) CreateReverseIPv6(ctx context.Context, instanceID string, reverseReq *govultr.ReverseIP) error {
	panic("implement me")
}

// ListReverseIPv6 is not implemented
func (f *FakeInstance) ListReverseIPv6(ctx context.Context, instanceID string) ([]govultr.ReverseIP, error) {
	panic("implement me")
}

// DeleteReverseIPv6 is not implemented
func (f *FakeInstance) DeleteReverseIPv6(ctx context.Context, instanceID, ip string) error {
	panic("implement me")
}

// CreateReverseIPv4 is not implemented
func (f *FakeInstance) CreateReverseIPv4(ctx context.Context, instanceID string, reverseReq *govultr.ReverseIP) error {
	panic("implement me")
}

// DefaultReverseIPv4 is not implemented
func (f *FakeInstance) DefaultReverseIPv4(ctx context.Context, instanceID, ip string) error {
	panic("implement me")
}

// GetUserData is not implemented
func (f *FakeInstance) GetUserData(ctx context.Context, instanceID string) (*govultr.UserData, error) {
	panic("implement me")
}

// GetUpgrades is not implemented
func (f *FakeInstance) GetUpgrades(ctx context.Context, instanceID string) (*govultr.Upgrades, error) {
	panic("implement me")
}
