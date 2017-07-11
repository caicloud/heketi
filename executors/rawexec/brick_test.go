//
// Copyright (c) 2016 The heketi Authors
//
// This file is licensed to you under your choice of the GNU Lesser
// General Public License, version 3 or any later version (LGPLv3 or
// later), or the GNU General Public License, version 2 (GPLv2), in all
// cases as published by the Free Software Foundation.
//

package caiexec

import (
	"strings"
	"testing"

	"github.com/heketi/heketi/executors"
	"github.com/heketi/heketi/pkg/utils"
	"github.com/heketi/tests"
)

func TestCaiExecBrickCreate(t *testing.T) {

	f := NewFakeSsh()
	defer tests.Patch(&sshNew,
		func(logger *utils.Logger, user string, file string) (Ssher, error) {
			return f, nil
		}).Restore()

	config := &CaiConfig{
		PrivateKeyFile: "xkeyfile",
		User:           "xuser",
		Port:           "100",
		CLICommandConfig: CLICommandConfig{
			Fstab: "/my/fstab",
		},
	}

	s, err := NewCaiExecutor(config)
	tests.Assert(t, err == nil)
	tests.Assert(t, s != nil)

	// Create a Brick
	b := &executors.BrickRequest{
		VgId:             "xvgid",
		Name:             "id",
		Device:           "/disk/sata01",
		TpSize:           100,
		Size:             10,
		PoolMetadataSize: 5,
	}

	// Mock ssh function
	f.FakeConnectAndExec = func(host string,
		commands []string,
		timeoutMinutes int,
		useSudo bool) ([]string, error) {

		tests.Assert(t, host == "myhost:100", host)
		tests.Assert(t, len(commands) == 1)

		for i, cmd := range commands {
			cmd = strings.Trim(cmd, " ")
			switch i {
			case 0:
				tests.Assert(t,
					cmd == "mkdir -p /disk/sata01/glusterfs/brick_id/brick", cmd)
			}
		}

		return nil, nil
	}

	// Create Brick
	_, err = s.BrickCreate("myhost", b)
	tests.Assert(t, err == nil, err)

}

func TestCaiExecBrickCreateWithGid(t *testing.T) {

	f := NewFakeSsh()
	defer tests.Patch(&sshNew,
		func(logger *utils.Logger, user string, file string) (Ssher, error) {
			return f, nil
		}).Restore()

	config := &CaiConfig{
		PrivateKeyFile: "xkeyfile",
		User:           "xuser",
		Port:           "100",
		CLICommandConfig: CLICommandConfig{
			Fstab: "/my/fstab",
		},
	}

	s, err := NewCaiExecutor(config)
	tests.Assert(t, err == nil)
	tests.Assert(t, s != nil)

	// Create a Brick
	b := &executors.BrickRequest{
		VgId:             "xvgid",
		Name:             "id",
		Device:           "/disk/sata01",
		TpSize:           100,
		Size:             10,
		PoolMetadataSize: 5,
		Gid:              1234,
	}

	// Mock ssh function
	f.FakeConnectAndExec = func(host string,
		commands []string,
		timeoutMinutes int,
		useSudo bool) ([]string, error) {

		tests.Assert(t, host == "myhost:100", host)
		tests.Assert(t, len(commands) == 3)

		for i, cmd := range commands {
			cmd = strings.Trim(cmd, " ")
			switch i {
			case 0:
				tests.Assert(t,
					cmd == "mkdir -p /disk/sata01/glusterfs/brick_id/brick", cmd)

			case 1:
				tests.Assert(t,
					cmd == "chown :1234 "+"/disk/sata01/glusterfs/brick_id/brick", cmd)

			case 2:
				tests.Assert(t,
					cmd == "chmod 2775 "+"/disk/sata01/glusterfs/brick_id/brick", cmd)
			}
		}

		return nil, nil
	}

	// Create Brick
	_, err = s.BrickCreate("myhost", b)
	tests.Assert(t, err == nil, err)

}

func TestCaiExecBrickCreateSudo(t *testing.T) {

	f := NewFakeSsh()
	defer tests.Patch(&sshNew,
		func(logger *utils.Logger, user string, file string) (Ssher, error) {
			return f, nil
		}).Restore()

	config := &CaiConfig{
		PrivateKeyFile: "xkeyfile",
		User:           "xuser",
		Port:           "100",
		CLICommandConfig: CLICommandConfig{
			Fstab: "/my/fstab",
			Sudo:  true,
		},
	}

	s, err := NewCaiExecutor(config)
	tests.Assert(t, err == nil)
	tests.Assert(t, s != nil)

	// Create a Brick
	b := &executors.BrickRequest{
		VgId:             "xvgid",
		Name:             "id",
		Device:           "/disk/sata01",
		TpSize:           100,
		Size:             10,
		PoolMetadataSize: 5,
	}

	// Mock ssh function
	f.FakeConnectAndExec = func(host string,
		commands []string,
		timeoutMinutes int,
		useSudo bool) ([]string, error) {

		tests.Assert(t, host == "myhost:100", host)
		tests.Assert(t, len(commands) == 1)
		tests.Assert(t, useSudo == true)

		for i, cmd := range commands {
			cmd = strings.Trim(cmd, " ")
			switch i {
			case 0:
				tests.Assert(t,
					cmd == "mkdir -p /disk/sata01/glusterfs/brick_id/brick", cmd)
			}
		}

		return nil, nil
	}

	// Create Brick
	_, err = s.BrickCreate("myhost", b)
	tests.Assert(t, err == nil, err)

}

func TestCaiExecBrickDestroy(t *testing.T) {

	f := NewFakeSsh()
	defer tests.Patch(&sshNew,
		func(logger *utils.Logger, user string, file string) (Ssher, error) {
			return f, nil
		}).Restore()

	config := &CaiConfig{
		PrivateKeyFile: "xkeyfile",
		User:           "xuser",
		Port:           "100",
		CLICommandConfig: CLICommandConfig{
			Fstab: "/my/fstab",
		},
	}

	s, err := NewCaiExecutor(config)
	tests.Assert(t, err == nil)
	tests.Assert(t, s != nil)

	// Create a Brick
	b := &executors.BrickRequest{
		VgId:             "xvgid",
		Name:             "id",
		Device:           "/disk/sata01",
		TpSize:           100,
		Size:             10,
		PoolMetadataSize: 5,
	}

	// Mock ssh function
	f.FakeConnectAndExec = func(host string,
		commands []string,
		timeoutMinutes int,
		useSudo bool) ([]string, error) {

		tests.Assert(t, host == "myhost:100", host)

		for _, cmd := range commands {
			cmd = strings.Trim(cmd, " ")
			switch {
			case strings.Contains(cmd, "rmdir"):
				tests.Assert(t,
					cmd == "rmdir /disk/sata01/glusterfs/brick_id/brick", cmd)
			}
		}

		return nil, nil
	}

	// Create Brick
	err = s.BrickDestroy("myhost", b)
	tests.Assert(t, err == nil, err)
}
