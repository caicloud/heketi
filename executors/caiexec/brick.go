//
// Copyright (c) 2015 The heketi Authors
//
// This file is licensed to you under your choice of the GNU Lesser
// General Public License, version 3 or any later version (LGPLv3 or
// later), or the GNU General Public License, version 2 (GPLv2), in all
// cases as published by the Free Software Foundation.
//

package caiexec

import (
	"fmt"

	"github.com/heketi/heketi/executors"
	"github.com/lpabon/godbc"
)

// Return the mount point for the brick
func (s *CaiExecutor) brickMountPoint(brick *executors.BrickRequest) string {
	return s.rootPath(brick.Device) + "/" + s.brickName(brick.Name) + "/brick"
}

// Device node for the lvm volume
func (s *CaiExecutor) devnode(brick *executors.BrickRequest) string {
	return s.rootPath(brick.Device) + "/" + s.brickName(brick.Name) + "/brick"
}

func (s *CaiExecutor) BrickCreate(host string,
	brick *executors.BrickRequest) (*executors.BrickInfo, error) {

	godbc.Require(brick != nil)
	godbc.Require(host != "")
	godbc.Require(brick.Name != "")
	godbc.Require(brick.Size > 0)
	// godbc.Require(brick.TpSize >= brick.Size)
	godbc.Require(brick.VgId != "")
	// godbc.Require(s.Fstab != "")

	// Create mountpoint name
	mountpoint := s.brickMountPoint(brick)

	// Create command set to execute on the node
	commands := []string{

		// Create a directory inside the formated volume for GlusterFS
		fmt.Sprintf("mkdir -p %v", mountpoint),
	}

	// Only set the GID if the value is other than root(gid 0).
	// When no gid is set, root is the only one that can write to the volume
	if 0 != brick.Gid {
		commands = append(commands, []string{
			// Set GID on brick
			fmt.Sprintf("chown :%v %v", brick.Gid, mountpoint),

			// Set writable by GID and UID
			fmt.Sprintf("chmod 2775 %v", mountpoint),
		}...)
	}

	// Execute commands
	_, err := s.RemoteExecutor.RemoteCommandExecute(host, commands, 10)
	if err != nil {
		// Cleanup
		s.BrickDestroy(host, brick)
		return nil, err
	}

	// Save brick location
	b := &executors.BrickInfo{
		Path: mountpoint,
	}
	return b, nil
}

func (s *CaiExecutor) BrickDestroy(host string,
	brick *executors.BrickRequest) error {

	godbc.Require(brick != nil)
	godbc.Require(host != "")
	godbc.Require(brick.Name != "")
	godbc.Require(brick.VgId != "")

	// Cleanup the mount point
	commands := []string{
		fmt.Sprintf("rmdir %v", s.brickMountPoint(brick)),
	}
	_, err := s.RemoteExecutor.RemoteCommandExecute(host, commands, 5)
	if err != nil {
		logger.Err(err)
	}

	return nil
}

func (s *CaiExecutor) BrickDestroyCheck(host string,
	brick *executors.BrickRequest) error {
	godbc.Require(brick != nil)
	godbc.Require(host != "")
	godbc.Require(brick.Name != "")
	godbc.Require(brick.VgId != "")

	return nil
}
