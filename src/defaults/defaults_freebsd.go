//go:build freebsd
// +build freebsd

package defaults

// Sane defaults for the BSD platforms. The "default" options may be
// may be replaced by the running configuration.
func getDefaults() platformDefaultParameters {
	return platformDefaultParameters{
		// Admin
		DefaultAdminListen: "tcp://localhost:9001",

		// Configuration (used for meshctl)
		DefaultConfigFile: "/usr/local/etc/mesh.conf",

		// Multicast interfaces
		DefaultMulticastInterfaces: []MulticastInterfaceConfig{
			{Regex: ".*", Beacon: true, Listen: true},
		},

                // Network domain
                DefaultNetworkDomain: NetworkDomainConfig{
                        Prefix: [...]byte{0xfc},
                },

		// TUN
		MaximumIfMTU:  32767,
		DefaultIfMTU:  32767,
		DefaultIfName: "/dev/tun0",
	}
}
