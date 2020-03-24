package driver

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

type Mounter interface {
	Format(source, fs string) error
	IsFormatted() (bool, error)
	Mount(source, target, fs string, opts ...string) error
	IsMounted() (bool, error)
	UnMount(target string) error
}

type mounter struct {
	log *logrus.Entry
}

func NewMounter(log *logrus.Entry) *mounter {
	return &mounter{log: log}
}

func (m *mounter) Format(source, fs string) error {
	if fs == "" {
		return errors.New("fs type was not provided - required for formatting the volume")
	}

	if source == "" {
		return errors.New("source type was not provided - required for formatting the volume")
	}

	mkFs := fmt.Sprintf("mkfs.%s", fs)
	_, err := exec.LookPath(mkFs)
	if err != nil {
		if err == exec.ErrNotFound {
			return fmt.Errorf("%q executable not found in $PATH", mkFs)
		}
		return err
	}

	argument := []string{}
	argument = append(argument, source)
	if fs == "ext4" || fs == "ext3" {
		argument = []string{"-F", source}
	}

	out, err := exec.Command(mkFs, argument...).CombinedOutput()

	if err != nil {
		return fmt.Errorf("formatting disk failed: %v cmd: '%s %s' output: %q",
			err, mkFs, strings.Join(argument, " "), string(out))
	}

	return nil
}

func (m *mounter) IsFormatted(target string) (bool, error) {
	if target == "" {
		return false, errors.New("source name was not provided")
	}

	blkidCmd := "blkid"
	_, err := exec.LookPath(blkidCmd)
	if err != nil {
		return false, fmt.Errorf("%q not found in $PATH", blkidCmd)
	}

	blkidArgs := []string{target}

	_, err = exec.Command(blkidCmd, blkidArgs...).Output()
	if err != nil {
		return false, fmt.Errorf("checking formatting failed for %v: %v", blkidArgs, err)
	}

	return true, nil
}

func (m *mounter) Mount(source, target, fs string, opts ...string) error {
	if source == "" {
		return errors.New("source type was not provided - required for mounting")
	}

	if target == "" {
		return errors.New("target type was not provided - required for mounting")
	}

	if fs == "" {
		return errors.New("fs type was not provided - required for mounting")
	}

	m.log.WithFields(logrus.Fields{
		"source":     source,
		"target":     target,
		"filesystem": fs,
		"options":    opts,
		"methods":    "mount",
	}).Info("Mount Called")

	mountCommand := "mount"
	mountArguments := []string{}

	mountArguments = append(mountArguments, "-t", fs)
	err := os.MkdirAll(target, 0750)
	if err != nil {
		return err
	}

	if len(opts) > 0 {
		mountArguments = append(mountArguments, "-o", strings.Join(opts, ","))
	}

	mountArguments = append(mountArguments, source)
	mountArguments = append(mountArguments, target)

	out, err := exec.Command(mountCommand, mountArguments...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("mounting failed: %v cmd: '%s %s' output: %q",
			err, mountCommand, strings.Join(mountArguments, " "), string(out))
	}

	return nil
}

func (m *mounter) IsMounted(path string) (bool, error) {
	if path == "" {
		return false, fmt.Errorf("No executable found in $PATH %v", path)
	}

	findmntCmd := "findmnt"
	_, err := exec.LookPath(findmntCmd)
	if err != nil {
		if err == exec.ErrNotFound {
			return false, fmt.Errorf("%q not found in $PATH", findmntCmd)
		}
		return false, err
	}

	cmdArgs := []string{"-o", "TARGET,PROPAGATION,FSTYPE,OPTIONS", "-M", path}
	out, err := exec.Command(path, cmdArgs...).CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("checking mount failed: ", err)
	}
	strOut := strings.Split(string(out), " ")[0]
	strOut = strings.TrimSuffix(string(out), "\n")

	if strOut == path {
		return true, nil
	}

	return false, nil
}

func (m *mounter) UnMount(target string) error {
	umountCmd := "umount"
	if target == "" {
		return errors.New("target is not specified for unmounting the volume")
	}

	umountArgs := []string{target}

	m.log.WithFields(logrus.Fields{
		"cmd":  umountCmd,
		"args": umountArgs,
	}).Info("executing umount command")

	out, err := exec.Command(umountCmd, umountArgs...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("unmounting failed: %v cmd: '%s %s' output: %q",
			err, umountCmd, target, string(out))
	}

	return nil
}
