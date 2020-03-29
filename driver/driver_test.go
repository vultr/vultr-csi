package driver

import (
	"context"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/sirupsen/logrus"
	"github.com/vultr/govultr"
	"golang.org/x/sync/errgroup"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestDriverSuite(t *testing.T) {
	socket := "/tmp/csi.sock"
	endpoint := "unix://" + socket
	if err := os.Remove(socket); err != nil && !os.IsNotExist(err) {
		t.Fatalf("failed to remove unix domain socket %s, error: %s", socket, err)
	}

	nodeID := "123456"
	region := "1"
	token := "dummy"
	version := "dev"
	client := govultr.NewClient(nil, token)

	log := logrus.New().WithFields(logrus.Fields{
		"region":  "1",
		"host_id": "12345",
		"version": "dev",
	})

	d := &VultrDriver{
		name:     DefaultDriverName,
		version:  version,
		endpoint: endpoint,

		client: client,
		nodeID: nodeID,
		region: region,

		waitTimeout: defaultTimeout,

		log:     log,
		mounter: NewFakeMounter(log),
	}

	go d.Run()

	_, cancel := context.WithCancel(context.Background())

	var eg errgroup.Group

	cancel()
	if err := eg.Wait(); err != nil {
		t.Errorf("driver run failed: %s", err)
	}
}

type fakeMounter struct {
	log     *logrus.Entry
	mounted map[string]string
}

type fakeStorageDriver struct {
	volumes map[string]*govultr.BlockStorage
}

func NewFakeMounter(log *logrus.Entry) *fakeMounter {
	return &fakeMounter{log: log}
}

func (f *fakeMounter) Format(source, fs string) error {
	return nil
}

func (f *fakeMounter) IsFormatted(source string) (bool, error) {
	return true, nil
}

func (f *fakeMounter) Mount(source, target, fs string, opts ...string) error {
	return nil
}

func (f *fakeMounter) IsMounted(target string) (bool, error) {
	return true, nil
}

func (f *fakeMounter) UnMount(target string) error {
	delete(f.mounted, target)
	return nil
}

func (f *fakeStorageDriver) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*govultr.BlockStorage, error) {
	volName := req.Name

	var curVolume *govultr.BlockStorage
	for _, volume := range f.volumes {
		if volume.Label == volName {
			curVolume = volume
		}
	}

	if curVolume != nil {
		return curVolume, nil
	}

	id := "123456"
	vol := &govultr.BlockStorage{
		BlockStorageID: id,
		RegionID:       1,
		Label:          volName,
		SizeGB:         10,
	}

	storage := make(map[string]*govultr.BlockStorage)
	f.volumes = storage
	f.volumes[id] = vol

	return vol, nil
}

func (f *fakeStorageDriver) DeleteVolume(ctx context.Context, id string) error {
	delete(f.volumes, id)
	return nil
}
