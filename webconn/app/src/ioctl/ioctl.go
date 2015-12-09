/**
 * The MIT License (MIT)
 *
 * Copyright (c) 2015 Edward Kim <edward@webconn.me>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package ioctl

import (
	"syscall"
)

const (
	IOC_NRBITS   = 8
	IOC_TYPEBITS = 8

	IOC_SIZEBITS = 14
	IOC_DIRBITS  = 2

	IOC_NRMASK   = (1 << IOC_NRBITS) - 1
	IOC_TYPEMASK = (1 << IOC_TYPEBITS) - 1
	IOC_SIZEMASK = (1 << IOC_SIZEBITS) - 1
	IOC_DIRMASK  = (1 << IOC_DIRBITS) - 1

	IOC_NRSHIFT   = 0
	IOC_TYPESHIFT = IOC_NRSHIFT + IOC_NRBITS
	IOC_SIZESHIFT = IOC_TYPESHIFT + IOC_TYPEBITS
	IOC_DIRSHIFT  = IOC_SIZESHIFT + IOC_SIZEBITS

	// Direction bits
	IOC_NONE  = 0
	IOC_WRITE = 1
	IOC_READ  = 2
)

//...and for the drivers/sound files...
const (
	IOC_IN        = IOC_WRITE << IOC_DIRSHIFT
	IOC_OUT       = IOC_READ << IOC_DIRSHIFT
	IOC_INOUT     = (IOC_WRITE | IOC_READ) << IOC_DIRSHIFT
	IOCSIZE_MASK  = IOC_SIZEMASK << IOC_SIZESHIFT
	IOCSIZE_SHIFT = IOC_SIZESHIFT
)

func IOC(dir, t, nr, size uintptr) uintptr {
	return (dir << IOC_DIRSHIFT) | (t << IOC_TYPESHIFT) | (nr << IOC_NRSHIFT) | (size << IOC_SIZESHIFT)
}

// used to create ioctl numbers

func IO(t, nr uintptr) uintptr {
	return IOC(IOC_NONE, t, nr, 0)
}

func IOR(t, nr, size uintptr) uintptr {
	return IOC(IOC_READ, t, nr, size)
}

func IOW(t, nr, size uintptr) uintptr {
	return IOC(IOC_WRITE, t, nr, size)
}

func IOWR(t, nr, size uintptr) uintptr {
	return IOC(IOC_READ|IOC_WRITE, t, nr, size)
}

func IOR_BAD(t, nr, size uintptr) uintptr {
	return IOC(IOC_READ, t, nr, size)
}

func IOW_BAD(t, nr, size uintptr) uintptr {
	return IOC(IOC_WRITE, t, nr, size)
}

func IOWR_BAD(t, nr, size uintptr) uintptr {
	return IOC(IOC_READ|IOC_WRITE, t, nr, size)
}

func IOCTL(fd, op, arg uintptr) error {
	_, _, ep := syscall.Syscall(syscall.SYS_IOCTL, fd, op, arg)
	if ep != 0 {
		return syscall.Errno(ep)
	}
	return nil
}