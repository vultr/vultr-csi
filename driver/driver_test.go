package driver

import (
	"context"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/kubernetes-csi/csi-test/pkg/sanity"
	"github.com/sirupsen/logrus"
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
	version := "dev"
	client := newFakeClient()

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
	randString := randString(5)

	cfg := &sanity.Config{
		TargetPath:  os.TempDir() + "/csi-target-" + randString,
		StagingPath: os.TempDir() + "/csi-staging-" + randString,
		Address:     endpoint,
	}
	sanity.Test(t, cfg)

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

func randString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
