// Copyright 2023 The NATS Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build !windows
// +build !windows

package server

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
)

// Acquire filesystem-level lock for the filestore to ensure exclusive access
// for this server instance.
func (fs *fileStore) lockFileSystem() error {
	var err error
	lpath := filepath.Join(fs.fcfg.StoreDir, "LOCK")
	fs.lfd, err = os.Create(lpath);
	if err != nil {
		return fmt.Errorf("could not create `%s': %v", lpath, err)
	}

	err = syscall.FcntlFlock(fs.lfd.Fd(), syscall.F_SETLK, &syscall.Flock_t{
		Type:    syscall.F_WRLCK,
		Whence:  io.SeekStart,
		Start:   0,
		Len:     0,
	})
	if err != nil {
		return fmt.Errorf("lock `%s': %v", fs.lfd.Name(), err)
	}

	return nil
}
