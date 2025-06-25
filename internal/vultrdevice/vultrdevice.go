// Package vultrdevice is used on a cluster node to ensure that the vultr
// storage devices are present and available in the configuration which is
// required by the CSI.
package vultrdevice

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"

	"github.com/vultr/vultr-csi/internal/vultruserdata"
)

const (
	symlinkPath   = "/dev/disk/by-id/virtio-"
	devRoot       = "/dev/"
	sysPCIPath    = "/sys/devices/pci0000:00/"
	sysSerialName = "serial"
)

type device struct {
	Name   string
	Serial string
}

// LinkBySerial iterates over all devices on the OS and if the serial matches
// what is provided, checks for a symlink to the device. it then creates a
// symlink if it does not already exist
func LinkBySerial(serial string) error {
	if runtime.GOOS != "linux" || !vultruserdata.IsVKE() {
		// the serial check is not relevant to this node
		return nil
	}

	devices, err := listSysDevices()
	if err != nil {
		return fmt.Errorf("unable to list and verify sys device info : %s", err)
	}

	for i := range devices {
		if devices[i].Serial == serial {
			// check if symlinked
			if _, err := os.Stat(fmt.Sprintf("%s%s", symlinkPath, serial)); err != nil {
				if os.IsNotExist(err) {
					// symlink doesn't exist; create it
					if err := os.Symlink(devices[i].Name, fmt.Sprintf("%s%s", symlinkPath, devices[i].Serial)); err != nil {
						return fmt.Errorf("unable to create symlink : %s", err)
					}

					return nil

				} else {
					// some other error, abort
					return fmt.Errorf("unable to read symlink : %s", err)
				}
			}

			// symlink exists; nothing to do
			return nil
		}
	}

	return fmt.Errorf("serial not found")
}

// listSysDevices traverses files in the /sys/devices/ directories, looks for any
// device that has a `serial` file and builds out the `device` struct with data
// used in the symlink, returning all matching devices
func listSysDevices() ([]device, error) {
	var devices []device
	var listSysSerial = func(path string, dirInfo fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// found a serial file, check it...
		if !dirInfo.IsDir() && dirInfo.Name() == sysSerialName {

			// this assumes that the serial file parent directory is the device
			// name and formats it under that assumption
			devName := fmt.Sprintf("%s%s", devRoot, filepath.Base(filepath.Dir(path)))

			// read what is in the 'serial' file; set that serial of the device
			devSerial, err := readSerial(path)
			if err != nil {
				return fmt.Errorf("unable to read serial file %q : %s", path, err)
			}

			// no serial, we don't care about this device
			if devSerial == "" {
				return nil
			}

			devices = append(devices, device{
				Name:   devName,
				Serial: devSerial,
			})
		}

		return nil
	}

	if err := filepath.WalkDir(sysPCIPath, listSysSerial); err != nil {
		return nil, fmt.Errorf("error walking sys dir files : %s", err)
	}

	return devices, nil
}

// readSerial is used to read the serial formatted file used in the /sys/devices
// directories. it will only return the first line of the file.
func readSerial(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("error reading file %q : %s", path, err)
	}

	defer f.Close()

	var serial string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		serial = scanner.Text()
		break
	}

	return serial, nil
}
