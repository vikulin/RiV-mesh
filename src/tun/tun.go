package tun

// This manages the tun driver to send/recv packets to/from applications

// TODO: Connection timeouts (call Conn.Close() when we want to time out)
// TODO: Don't block in reader on writes that are pending searches

import (
	"errors"
	"fmt"
	"net"

	"github.com/Arceliar/phony"
	"golang.zx2c4.com/wireguard/tun"

	//"github.com/RiV-chain/RiV-mesh/src/address"
	"github.com/RiV-chain/RiV-mesh/src/core"
	"github.com/RiV-chain/RiV-mesh/src/defaults"
	"github.com/RiV-chain/RiV-mesh/src/ipv6rwc"
)

type MTU uint16

// TunAdapter represents a running TUN interface and extends the
// mesh.Adapter type. In order to use the TUN adapter with Mesh, you
// should pass this object to the mesh.SetRouterAdapter() function before
// calling mesh.Start().
type TunAdapter struct {
	core        *core.Core
	rwc         *ipv6rwc.ReadWriteCloser
	log         core.Logger
	addr        core.Address
	subnet      core.Subnet
	mtu         uint64
	iface       tun.Device
	phony.Inbox // Currently only used for _handlePacket from the reader, TODO: all the stuff that currently needs a mutex below
	//mutex        sync.RWMutex // Protects the below
	isOpen    bool
	isEnabled bool // Used by the writer to drop sessionTraffic if not enabled
	config    struct {
		name InterfaceName
		mtu  InterfaceMTU
	}
}

// Gets the maximum supported MTU for the platform based on the defaults in
// defaults.GetDefaults().
func getSupportedMTU(mtu uint64) uint64 {
	if mtu < 1280 {
		return 1280
	}
	if mtu > MaximumMTU() {
		return MaximumMTU()
	}
	return mtu
}

// Name returns the name of the adapter, e.g. "tun0". On Windows, this may
// return a canonical adapter name instead.
func (tun *TunAdapter) Name() string {
	if name, err := tun.iface.Name(); err == nil {
		return name
	}
	return ""
}

// MTU gets the adapter's MTU. This can range between 1280 and 65535, although
// the maximum value is determined by your platform. The returned value will
// never exceed that of MaximumMTU().
func (tun *TunAdapter) MTU() uint64 {
	return getSupportedMTU(tun.mtu)
}

// DefaultName gets the default TUN interface name for your platform.
func DefaultName() string {
	return defaults.GetDefaults().DefaultIfName
}

// DefaultMTU gets the default TUN interface MTU for your platform. This can
// be as high as MaximumMTU(), depending on platform, but is never lower than 1280.
func DefaultMTU() uint64 {
	return defaults.GetDefaults().DefaultIfMTU
}

// MaximumMTU returns the maximum supported TUN interface MTU for your
// platform. This can be as high as 65535, depending on platform, but is never
// lower than 1280.
func MaximumMTU() uint64 {
	return defaults.GetDefaults().MaximumIfMTU
}

// Init initialises the TUN module. You must have acquired a Listener from
// the Mesh	 core before this point and it must not be in use elsewhere.
func New(core *core.Core, log core.Logger, opts ...SetupOption) (*TunAdapter, error) {
	tun := &TunAdapter{
		core: core,
		rwc: ipv6rwc.NewReadWriteCloser(core),
		log: log,
	}
	for _, opt := range opts {
		tun._applyOption(opt)
	}
	return tun, tun._start()
}

func (tun *TunAdapter) _start() error {
	if tun.isOpen {
		return errors.New("TUN module is already started")
	}
	tun.addr = tun.rwc.Address()
	tun.subnet = tun.rwc.Subnet()
	addr := fmt.Sprintf("%s/%d", net.IP(tun.addr[:]).String(), 8*len(tun.core.GetPrefix())-1)
	if tun.config.name == "none" || tun.config.name == "dummy" {
		tun.log.Debugln("Not starting TUN as ifname is none or dummy")
		tun.isEnabled = false
		go tun.write()
		return nil
	}
	mtu := uint64(tun.config.mtu)
	if tun.rwc.MaxMTU() < mtu {
		mtu = tun.rwc.MaxMTU()
	}
	if err := tun.setup(string(tun.config.name), addr, mtu); err != nil {
		return err
	}
	if tun.MTU() != mtu {
		tun.log.Warnf("Warning: Interface MTU %d automatically adjusted to %d (supported range is 1280-%d)", tun.config.mtu, tun.MTU(), MaximumMTU())
	}
	tun.rwc.SetMTU(tun.MTU())
	tun.isOpen = true
	tun.isEnabled = true
	go tun.read()
	go tun.write()
	return nil
}

// IsStarted returns true if the module has been started.
func (tun *TunAdapter) IsStarted() bool {
	var isOpen bool
	phony.Block(tun, func() {
		isOpen = tun.isOpen
	})
	return isOpen
}

// Start the setup process for the TUN adapter. If successful, starts the
// read/write goroutines to handle packets on that interface.
func (tun *TunAdapter) Stop() error {
	var err error
	phony.Block(tun, func() {
		err = tun._stop()
	})
	return err
}

func (tun *TunAdapter) _stop() error {
	tun.isOpen = false
	// by TUN, e.g. readers/writers, sessions
	if tun.iface != nil {
		// Just in case we failed to start up the iface for some reason, this can apparently happen on Windows
		tun.iface.Close()
	}
	return nil
}
