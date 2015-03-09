// Copyright (c) 2012 VMware, Inc.

package sigar

// #include <stdlib.h>
// #include <windows.h>
import "C"

import (
	"bytes"
	"fmt"
	"unsafe"
)

func init() {
}

func (self *LoadAverage) Get() error {
	return nil
}

func (self *Uptime) Get() error {
	return nil
}

func (self *Mem) Get() error {
	var statex C.MEMORYSTATUSEX
	statex.dwLength = C.DWORD(unsafe.Sizeof(statex))

	succeeded := C.GlobalMemoryStatusEx(&statex)
	if succeeded == C.FALSE {
		lastError := C.GetLastError()
		return fmt.Errorf("GlobalMemoryStatusEx failed with error: %d", int(lastError))
	}

	self.Total = uint64(statex.ullTotalPhys)
	return nil
}

func (self *Swap) Get() error {
	return notImplemented()
}

func (self *Cpu) Get() error {
	return notImplemented()
}

func (self *CpuList) Get() error {
	return notImplemented()
}

func (self *FileSystemList) Get() error {
	capacity := len(self.List)
	if capacity == 0 {
		capacity = 10
	}
	fslist := make([]FileSystem, 0, capacity)

	buffer := C.CString("")
	defer C.free(unsafe.Pointer(buffer))

	buffer_length := C.GetLogicalDriveStrings(1024, (*C.CHAR)(buffer))
	if buffer_length == 0 {
		return fmt.Errorf("GetLogicalDriveStrings failed: %d", C.GetLastError())
	} else if buffer_length > 1024 {
		buffer_length = C.GetLogicalDriveStrings(buffer_length+1, (*C.CHAR)(buffer))
	}

	drives := C.GoBytes(unsafe.Pointer(buffer), (C.int)(buffer_length))
	drivename := new(bytes.Buffer)

	for _, ch := range drives {
		if ch != 0 {
			drivename.WriteByte(ch)
		} else {
			fs := FileSystem{}
			fs.DevName = drivename.String()
			fs.DirName = fs.DevName

			fslist = append(fslist, fs)

			drivename = new(bytes.Buffer)
		}
	}

	self.List = fslist

	return nil
}

func (self *ProcList) Get() error {
	return notImplemented()
}

func (self *ProcState) Get(pid int) error {
	return notImplemented()
}

func (self *ProcMem) Get(pid int) error {
	return notImplemented()
}

func (self *ProcTime) Get(pid int) error {
	return notImplemented()
}

func (self *ProcArgs) Get(pid int) error {
	return notImplemented()
}

func (self *ProcExe) Get(pid int) error {
	return notImplemented()
}

func (self *FileSystemUsage) Get(path string) error {
	var availableBytes C.ULARGE_INTEGER
	var totalBytes C.ULARGE_INTEGER
	var totalFreeBytes C.ULARGE_INTEGER

	pathChars := C.CString(path)
	defer C.free(unsafe.Pointer(pathChars))

	succeeded := C.GetDiskFreeSpaceEx((*C.CHAR)(pathChars), &availableBytes, &totalBytes, &totalFreeBytes)
	if succeeded == C.FALSE {
		lastError := C.GetLastError()
		return fmt.Errorf("GetDiskFreeSpaceEx failed with error: %d", int(lastError))
	}

	self.Total = *(*uint64)(unsafe.Pointer(&totalBytes))
	self.Free = *(*uint64)(unsafe.Pointer(&totalFreeBytes))
	self.Avail = *(*uint64)(unsafe.Pointer(&availableBytes))
	self.Used = self.Total - self.Free
	self.Files = 0
	self.FreeFiles = 0

	return nil
}

func notImplemented() error {
	panic("Not Implemented")
	return nil
}
