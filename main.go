package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"unsafe"

	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err.Error())
	}

	for {
		time.Sleep(time.Second)
	}
}

func run() error {
	fmt.Printf("Hello World!\n")

	// Before we can configure ethernet we need to load hardware drivers
	if err := addDriverModule(); err != nil {
		return errors.Wrap(err, "failed to add driver")
	}

	if err := configureEthernet(); err != nil {
		return errors.Wrap(err, "failed to configure ethernet")
	}

	fmt.Printf("Ethernet configured\n")

	http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from Scratch Machine!\n")
	}))

	return nil
}

var fakeString [3]byte

func addDriverModule() error {
	// We need a file descriptor for our file
	// driverPath := "/lib/modules/4.9.73-0-virthardened/kernel/drivers/net/ethernet/intel/e1000/e1000.ko"
	driverPath := "/e1000.ko"
	f, err := os.Open(driverPath)
	if err != nil {
		return errors.Wrap(err, "open of driver file failed")
	}
	defer f.Close()
	fd := f.Fd()

	_, _, errno := unix.Syscall(unix.SYS_FINIT_MODULE, fd, uintptr(unsafe.Pointer(&fakeString)), 0)
	if errno != 0 && errno != unix.EEXIST {
		return errors.Wrap(errno, "init module failed")
	}

	return nil
}

type socketAddrRequest struct {
	name [unix.IFNAMSIZ]byte
	addr unix.RawSockaddrInet4
}

type socketFlagsRequest struct {
	name  [unix.IFNAMSIZ]byte
	flags uint16
	pad   [22]byte
}

func configureEthernet() error {
	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		return errors.Wrap(err, "could not open control socket")
	}

	defer unix.Close(fd)

	// We want to associate an IP address with eth0, then set flags to
	// activate it

	sa := socketAddrRequest{}
	copy(sa.name[:], "eth0")
	sa.addr.Family = unix.AF_INET
	copy(sa.addr.Addr[:], []byte{10, 0, 2, 15})

	// Set address
	if err := ioctl(fd, unix.SIOCSIFADDR, uintptr(unsafe.Pointer(&sa))); err != nil {
		return errors.Wrap(err, "failed setting address for eth0")
	}

	// Set netmask
	copy(sa.addr.Addr[:], []byte{255, 255, 255, 0})
	if err := ioctl(fd, unix.SIOCSIFNETMASK, uintptr(unsafe.Pointer(&sa))); err != nil {
		return errors.Wrap(err, "failed setting netmask for eth0")
	}

	// Get flags
	sf := socketFlagsRequest{}
	sf.name = sa.name
	if err := ioctl(fd, unix.SIOCGIFFLAGS, uintptr(unsafe.Pointer(&sf))); err != nil {
		return errors.Wrap(err, "failed getting flags for eth0")
	}

	sf.flags |= unix.IFF_UP | unix.IFF_RUNNING
	if err := ioctl(fd, unix.SIOCSIFFLAGS, uintptr(unsafe.Pointer(&sf))); err != nil {
		return errors.Wrap(err, "failed getting flags for eth0")
	}

	return nil
}

func ioctl(fd int, code, data uintptr) error {
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), code, data)
	if errno != 0 {
		return errno
	}
	return nil
}
