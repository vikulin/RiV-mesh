//go:build windows
// +build windows

package tun

import (
	"bytes"
	"errors"
	"time"
	"log"
	"net"

	"github.com/RiV-chain/RiV-mesh/src/defaults"
	"golang.org/x/sys/windows"

	wgtun "golang.zx2c4.com/wireguard/tun"
	"golang.zx2c4.com/wireguard/windows/elevate"
	"golang.zx2c4.com/wireguard/windows/tunnel/winipcfg"
)

// This is to catch Windows platforms

// Configures the TUN adapter with the correct IPv6 address and MTU.
func (tun *TunAdapter) setup(ifname string, addr string, mtu uint64) error {
	if ifname == "auto" {
		ifname = defaults.GetDefaults().DefaultIfName
	}
	return elevate.DoAsSystem(func() error {
		var err error
		var iface wgtun.Device
		var guid windows.GUID
		if guid, err = windows.GUIDFromString("{f1369c05-0344-40ed-a772-bfb4770abdd0}"); err != nil {
			return err
		}
		if iface, err = wgtun.CreateTUNWithRequestedGUID(ifname, &guid, int(mtu)); err != nil {
			return err
		}
		tun.iface = iface
		for i := 1; i < 10; i++ {
			if err = tun.setupAddress(addr); err != nil {
				tun.log.Errorln("Failed to set up TUN address:", err)
				log.Printf("waiting...")
				if i > 8 {
					return err
				} else {
					time.Sleep(time.Duration(2 * i) * time.Second)
				}
			} else {
				break
			}
		}
		if err = tun.setupMTU(getSupportedMTU(mtu)); err != nil {
			tun.log.Errorln("Failed to set up TUN MTU:", err)
			return err
		}
		if mtu, err := iface.MTU(); err == nil {
			tun.mtu = uint64(mtu)
		}
		return nil
	})
}

// Sets the MTU of the TUN adapter.
func (tun *TunAdapter) setupMTU(mtu uint64) error {
	if tun.iface == nil || tun.Name() == "" {
		return errors.New("Can't configure MTU as TUN adapter is not present")
	}
	if intf, ok := tun.iface.(*wgtun.NativeTun); ok {
		luid := winipcfg.LUID(intf.LUID())
		ipfamily, err := luid.IPInterface(windows.AF_INET6)
		if err != nil {
			return err
		}

		ipfamily.NLMTU = uint32(mtu)
		intf.ForceMTU(int(ipfamily.NLMTU))
		ipfamily.UseAutomaticMetric = false
		ipfamily.Metric = 0
		ipfamily.DadTransmits = 0
		ipfamily.RouterDiscoveryBehavior = winipcfg.RouterDiscoveryDisabled

		if err := ipfamily.Set(); err != nil {
			return err
		}
	}

	return nil
}

// Sets the IPv6 address of the TUN adapter.
func (tun *TunAdapter) setupAddress(addr string) error {
	if tun.iface == nil || tun.Name() == "" {
		return errors.New("Can't configure IPv6 address as TUN adapter is not present")
	}
	if intf, ok := tun.iface.(*wgtun.NativeTun); ok {
		if ipaddr, ipnet, err := net.ParseCIDR(addr); err == nil {
			luid := winipcfg.LUID(intf.LUID())
			addresses := append([]net.IPNet{}, net.IPNet{
				IP:   ipaddr,
				Mask: ipnet.Mask,
			})

			err := luid.SetIPAddressesForFamily(windows.AF_INET6, addresses)
			if err == windows.ERROR_OBJECT_ALREADY_EXISTS {
				cleanupAddressesOnDisconnectedInterfaces(windows.AF_INET6, addresses)
				err = luid.SetIPAddressesForFamily(windows.AF_INET6, addresses)
			}
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		return errors.New("unable to get NativeTUN")
	}
	return nil
}

/*
 * cleanupAddressesOnDisconnectedInterfaces
 * SPDX-License-Identifier: MIT
 * Copyright (C) 2019 WireGuard LLC. All Rights Reserved.
 */
func cleanupAddressesOnDisconnectedInterfaces(family winipcfg.AddressFamily, addresses []net.IPNet) {
	if len(addresses) == 0 {
		return
	}
	includedInAddresses := func(a net.IPNet) bool {
		// TODO: this makes the whole algorithm O(n^2). But we can't stick net.IPNet in a Go hashmap. Bummer!
		for _, addr := range addresses {
			ip := addr.IP
			if ip4 := ip.To4(); ip4 != nil {
				ip = ip4
			}
			mA, _ := addr.Mask.Size()
			mB, _ := a.Mask.Size()
			if bytes.Equal(ip, a.IP) && mA == mB {
				return true
			}
		}
		return false
	}
	interfaces, err := winipcfg.GetAdaptersAddresses(family, winipcfg.GAAFlagDefault)
	if err != nil {
		return
	}
	for _, iface := range interfaces {
		if iface.OperStatus == winipcfg.IfOperStatusUp {
			continue
		}
		for address := iface.FirstUnicastAddress; address != nil; address = address.Next {
			ip := address.Address.IP()
			ipnet := net.IPNet{IP: ip, Mask: net.CIDRMask(int(address.OnLinkPrefixLength), 8*len(ip))}
			if includedInAddresses(ipnet) {
				log.Printf("Cleaning up stale address %s from interface ‘%s’", ipnet.String(), iface.FriendlyName())
				iface.LUID.DeleteIPAddress(ipnet)
			}
		}
	}
}
