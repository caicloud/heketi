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
)

// Assume the disk is formatted before using
func (s *CaiExecutor) DeviceSetup(host, rootPath, vgid string) (d *executors.DeviceInfo, e error) {

	// Setup commands
	commands := []string{
		fmt.Sprintf("mkdir -p /%s/%s", rootPath, vgid),
	}

	// Execute command
	_, err := s.RemoteExecutor.RemoteCommandExecute(host, commands, 5)
	if err != nil {
		return nil, err
	}

	// Create a cleanup function if anything fails
	defer func() {
		if e != nil {
			s.DeviceTeardown(host, rootPath, vgid)
		}
	}()

	d = &executors.DeviceInfo{}
	err = s.getDiskSizeFromNode(d, host, rootPath)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (s *CaiExecutor) DeviceTeardown(host, rootPath, vgid string) error {

	// Setup commands
	commands := []string{
		fmt.Sprintf("rm -rf /%s/%s", rootPath, vgid),
	}

	// Execute command
	_, err := s.RemoteExecutor.RemoteCommandExecute(host, commands, 5)
	if err != nil {
		logger.LogError("Error while deleting path %v on %v with id %v",
			rootPath, host, vgid)
	}

	return nil
}

func (s *CaiExecutor) getDiskSizeFromNode(d *executors.DeviceInfo, _host, _rootPath string) error {

	// TBD: No limit of disk size(1PB)
	d.Size = 1024 * 1024 * 1024 * 1024 * 1024
	d.ExtentSize = 4096
	return nil
}
