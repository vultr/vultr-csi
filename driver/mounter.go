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
	IsFormatted(source string) (bool, error)
	Mount(source, target, fs string, opts ...string) error
	IsMounted(target string) (bool, error)
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

	m.log.WithFields(logrus.Fields{
		"source":      source,
		"fs-type":     fs,
		"format-cmd":  mkFs,
		"format-args": argument,
	}).Info("Format called")

	out, err := exec.Command(mkFs, argument...).CombinedOutput()

	if err != nil {
		return fmt.Errorf("formatting disk failed: %v cmd: '%s %s' output: %q",
			err, mkFs, strings.Join(argument, " "), string(out))
	}

	return nil
}

func (m *mounter) IsFormatted(source string) (bool, error) {
	if source == "" {
		return false, errors.New("source name was not provided")
	}

	blkidCmd := "blkid"
	_, err := exec.LookPath(blkidCmd)
	if err != nil {
		return false, fmt.Errorf("%q not found in $PATH", blkidCmd)
	}

	blkidArgs := []string{source}

	m.log.WithFields(logrus.Fields{
		"format-command": blkidCmd,
		"format-args":    blkidArgs,
	}).Info("isFormatted called")

	out, err := exec.Command(blkidCmd, blkidArgs...).CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("checking formatting failed for %v: %v", blkidArgs, err)
	}

	// assume not formatted
	if string(out) == "" {
		return false, nil
	}

	m.log.WithFields(logrus.Fields{
		"format-output": out,
	}).Info("isFormatted end")
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

	m.log.WithFields(logrus.Fields{
		"mount command":   mountCommand,
		"mount arguments": mountArguments,
	}).Info("mount command and arguments")

	out, err := exec.Command(mountCommand, mountArguments...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("mounting failed: %v cmd: '%s %s' output: %q",
			err, mountCommand, strings.Join(mountArguments, " "), string(out))
	}

	if _, err := os.Stat(target + "/lost+found"); err == nil {
		os.Remove(target + "/lost+found")
	} else if os.IsNotExist(err) {
		m.log.WithFields(logrus.Fields{
			"error": err,
		}).Info("mount command - removal of lost+found error")
	}

	return nil
}

func (m *mounter) IsMounted(target string) (bool, error) {
	if target == "" {
		return false, errors.New("target path was not provided")
	}

	findmntCmd := "findmnt"
	_, err := exec.LookPath(findmntCmd)
	if err != nil {
		if err == exec.ErrNotFound {
			return false, fmt.Errorf("%q not found in $PATH", findmntCmd)
		}
		return false, err
	}

	cmdArgs := []string{"-o", "TARGET", "-T", target}
	out, err := exec.Command(findmntCmd, cmdArgs...).CombinedOutput()
	if err != nil {
		// not an error, just nothing found.
		if strings.TrimSpace(string(out)) == "" {
			return false, nil
		}

		return false, fmt.Errorf("checking mount failed with command %v: %v", findmntCmd, err)
	}

	if string(out) == "" {
		return false, nil
	}

	if strings.Contains(string(out), target) {
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
