package driver

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
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

func (m *mounter) IsFormatted() (bool, error) {
	panic("implement me")
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

func (m *mounter) IsMounted() (bool, error) {
	panic("implement me")
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
