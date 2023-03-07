package driver

import (
	"context"
	"math/rand"
	"os"
	"testing"
	"time"

	"golang.org/x/oauth2"

	"github.com/sirupsen/logrus"
	"github.com/vultr/govultr/v2"
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

	nodeID := "245bb2fe-b55c-44a0-9a1e-ab80e4b5f088"
	region := "ewr"
	token := "dummy"
	version := "dev"
	ctx := context.Background()
	config := &oauth2.Config{}
	ts := config.TokenSource(ctx, &oauth2.Token{AccessToken: token})
	client := govultr.NewClient(oauth2.NewClient(ctx, ts))

	log := logrus.New().WithFields(logrus.Fields{
		"region":  "ewr",
		"host_id": "245bb2fe-b55c-44a0-9a1e-ab80e4b5f088",
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

func (f *fakeMounter) GetStatistics(volumePath string) (volumeStatistics, error) {
	return volumeStatistics{
		availableBytes: 3 * giB,
		totalBytes:     10 * giB,
		usedBytes:      7 * giB,

		availableInodes: 3000,
		totalInodes:     10000,
		usedInodes:      7000,
	}, nil
}

func (f *fakeMounter) IsBlockDevice(volumePath string) (bool, error) {
	return false, nil
}
