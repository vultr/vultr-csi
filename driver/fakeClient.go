package driver

import (
	"context"
	"net/http"

	"github.com/vultr/govultr/v3"
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

func (f *fakeBS) Create(ctx context.Context, blockReq *govultr.BlockStorageCreate) (*govultr.BlockStorage, *http.Response, error) {
	return newFakeBS(), nil, nil
}

func (f *fakeBS) Get(ctx context.Context, blockID string) (*govultr.BlockStorage, *http.Response, error) {
	return newFakeBS(), nil, nil
}

func (f *fakeBS) Update(ctx context.Context, blockID string, blockReq *govultr.BlockStorageUpdate) error {
	panic("implement me")
}

func (f *fakeBS) Delete(ctx context.Context, blockID string) error {
	return nil
}

func (f *fakeBS) List(ctx context.Context, options *govultr.ListOptions) ([]govultr.BlockStorage, *govultr.Meta, *http.Response, error) {
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
		}, nil, nil
}

func (f *fakeBS) Attach(ctx context.Context, blockID string, attach *govultr.BlockStorageAttach) error {
	panic("implement me")
}

func (f *fakeBS) Detach(ctx context.Context, blockID string, detach *govultr.BlockStorageDetach) error {
	list, _, _, err := f.List(ctx, nil) //nolint:bodyclose
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
func (f *FakeInstance) Create(ctx context.Context, instanceReq *govultr.InstanceCreateReq) (*govultr.Instance, *http.Response, error) {
	panic("implement me")
}

// Get returns an instance struct
func (f *FakeInstance) Get(ctx context.Context, instanceID string) (*govultr.Instance, *http.Response, error) {
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
	}, nil, nil
}

// Update updates and instance
func (f *FakeInstance) Update(_ context.Context, _ string, _ *govultr.InstanceUpdateReq) (*govultr.Instance, *http.Response, error) {
	panic("implement me")
}

// Delete jis not implemented
func (f *FakeInstance) Delete(ctx context.Context, instanceID string) error {
	panic("implement me")
}

// List is not implemented
func (f *FakeInstance) List(ctx context.Context, options *govultr.ListOptions) ([]govultr.Instance, *govultr.Meta, *http.Response, error) {
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

// Reinstall reinstalls an instance
func (f *FakeInstance) Reinstall(_ context.Context, _ string, _ *govultr.ReinstallReq) (*govultr.Instance, *http.Response, error) {
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

// Restore restores an instance
func (f *FakeInstance) Restore(_ context.Context, _ string, _ *govultr.RestoreReq) (*http.Response, error) {
	panic("implement me")
}

// GetBandwidth gets bandwidth for an instance
func (f *FakeInstance) GetBandwidth(_ context.Context, _ string) (*govultr.Bandwidth, *http.Response, error) {
	panic("implement me")
}

// GetNeighbors gets neighors for an instance
func (f *FakeInstance) GetNeighbors(_ context.Context, _ string) (*govultr.Neighbors, *http.Response, error) {
	panic("implement me")
}

// ListVPCInfo is not implemented
func (f *FakeInstance) ListVPCInfo(ctx context.Context, instanceID string, options *govultr.ListOptions) ([]govultr.VPCInfo, *govultr.Meta, *http.Response, error) { //nolint:lll
	panic("implement me")
}

// AttachVPC is not implemented
func (f *FakeInstance) AttachVPC(ctx context.Context, instanceID, networkID string) error {
	panic("implement me")
}

// DetachVPC is not implemented
func (f *FakeInstance) DetachVPC(ctx context.Context, instanceID, networkID string) error {
	panic("implement me")
}

// ListVPC2Info is not implemented
func (f *FakeInstance) ListVPC2Info(ctx context.Context, instanceID string, options *govultr.ListOptions) ([]govultr.VPC2Info, *govultr.Meta, *http.Response, error) { //nolint:lll
	panic("implement me")
}

// AttachVPC2 is not implemented
func (f *FakeInstance) AttachVPC2(ctx context.Context, instanceID string, vpc2Req *govultr.AttachVPC2Req) error {
	panic("implement me")
}

// DetachVPC2 is not implemented
func (f *FakeInstance) DetachVPC2(ctx context.Context, instanceID, vpcID string) error {
	panic("implement me")
}

// ISOStatus gets ISO status from instance
func (f *FakeInstance) ISOStatus(_ context.Context, _ string) (*govultr.Iso, *http.Response, error) {
	panic("implement me")
}

// AttachISO attaches ISO to instance
func (f *FakeInstance) AttachISO(_ context.Context, _, _ string) (*http.Response, error) {
	panic("implement me")
}

// DetachISO detaches ISO from instance
func (f *FakeInstance) DetachISO(_ context.Context, _ string) (*http.Response, error) {
	panic("implement me")
}

// GetBackupSchedule gets instance backup stchedule
func (f *FakeInstance) GetBackupSchedule(_ context.Context, _ string) (*govultr.BackupSchedule, *http.Response, error) {
	panic("implement me")
}

// SetBackupSchedule sets instance backup schedule
func (f *FakeInstance) SetBackupSchedule(_ context.Context, _ string, _ *govultr.BackupScheduleReq) (*http.Response, error) {
	panic("implement me")
}

// CreateIPv4 creates an IPv4 association to instance
func (f *FakeInstance) CreateIPv4(_ context.Context, _ string, _ *bool) (*govultr.IPv4, *http.Response, error) {
	panic("implement me")
}

// ListIPv4 gets IPv4 addresses associated with instance
func (f *FakeInstance) ListIPv4(_ context.Context, _ string, _ *govultr.ListOptions) ([]govultr.IPv4, *govultr.Meta, *http.Response, error) { //nolint:lll
	panic("implement me")
}

// DeleteIPv4 is not implemented
func (f *FakeInstance) DeleteIPv4(ctx context.Context, instanceID, ip string) error {
	panic("implement me")
}

// ListIPv6 lists IPv6 addresses associated with instance
func (f *FakeInstance) ListIPv6(_ context.Context, _ string, _ *govultr.ListOptions) ([]govultr.IPv6, *govultr.Meta, *http.Response, error) { //nolint:lll
	panic("implement me")
}

// CreateReverseIPv6 is not implemented
func (f *FakeInstance) CreateReverseIPv6(ctx context.Context, instanceID string, reverseReq *govultr.ReverseIP) error {
	panic("implement me")
}

// ListReverseIPv6 gets reverse IP for IPv6 on instance
func (f *FakeInstance) ListReverseIPv6(_ context.Context, _ string) ([]govultr.ReverseIP, *http.Response, error) {
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

// GetUserData returns instance userdata
func (f *FakeInstance) GetUserData(_ context.Context, _ string) (*govultr.UserData, *http.Response, error) {
	panic("implement me")
}

// GetUpgrades gets instance upgade
func (f *FakeInstance) GetUpgrades(_ context.Context, _ string) (*govultr.Upgrades, *http.Response, error) {
	panic("implement me")
}
