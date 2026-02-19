# RiV-mesh - Decentralized IPv6 Mesh Network

## Introduction

RiV-mesh is the legacy implementation of a fully end-to-end encrypted IPv6
mesh network, designed to provide secure connectivity between a wide spectrum of endpoint devices like IoT devices,
desktop computers or even routers.
It is lightweight, self-arranging, supported on multiple
platforms and allows pretty much any IPv6-capable application
to communicate securely with other network nodes.

### RiV-mesh
- **Status**: Legacy maintenance mode
- **License**: LGPLv3 only
- **Features**: Core mesh networking functionality
- **Use Case**: Existing deployments, compatibility

## Core Functionality

RiV-mesh provides the decentralized IPv6 mesh overlay network infrastructure for applications like CupLink™, which requires a VPN service to establish the mesh network environment for peer-to-peer communication.

## Supported Platforms

RiV-mesh works on a number of platforms, including Linux, macOS, Ubiquiti
EdgeRouter, VyOS, Windows, FreeBSD, OpenBSD and OpenWrt.

Please see our [Installation](https://github.com/vikulin/RiV-mesh-builds#riv-mesh-build) 
page for more information. You may also find other platform-specific wrappers, scripts
or tools in the `contrib` folder.

## Building

If you want to build from source, as opposed to installing one of the pre-built
packages:

1. Install [Go](https://golang.org) (requires Go 1.19 or later)
2. Clone this repository
2. Run `./build`

Note that you can cross-compile for other platforms and architectures by
specifying the `GOOS` and `GOARCH` environment variables, e.g. `GOOS=windows
./build` or `GOOS=linux GOARCH=mipsle ./build`

... or generate an iOS framework with:

```
./contrib/mobile/build -i
```

... or generate an Android AAR bundle with:

```
./contrib/mobile/build -a
```

Other OS packages can be built in this repo: https://github.com/vikulin/RiV-mesh-builds.

## Running

### Generate configuration

To generate static configuration, either generate a HJSON file (human-friendly,
complete with comments):

```
./mesh -genconf > /path/to/mesh.conf
```

... or generate a plain JSON file (which is easy to manipulate
programmatically):

```
./mesh -genconf -json > /path/to/mesh.conf
```

You will need to edit the `mesh.conf` file to add or remove peers, modify
other configuration such as listen addresses or multicast addresses, etc.

### Run RiV-mesh

To run with the generated static configuration:

```
./mesh -useconffile /path/to/mesh.conf
```

To run in auto-configuration mode (which will use sane defaults and random keys
at each startup, instead of using a static configuration file):

```
./mesh -autoconf
```

You will likely need to run RiV-mesh as a privileged user or under `sudo`,
unless you have permission to create TUN/TAP adapters. On Linux this can be done
by giving the RiV-mesh binary the `CAP_NET_ADMIN` capability.

## Documentation

Documentation is available [on our website](https://vikulin.github.io/RiV-mesh/).

- [Installing RiV-mesh](https://vikulin.github.io/RiV-mesh/)
- [Configuring RiV-mesh](https://vikulin.github.io/RiV-mesh/)
- [Frequently asked questions](https://vikulin.github.io/RiV-mesh/)
- [Version changelog](CHANGELOG.md)

## Work in progress:

<img width="253" alt="DDNS" src="https://github.com/vikulin/RiV-mesh/assets/743622/05757b9f-2053-4503-9037-24e5ff992554">

## Public peers
If you are operating a RiV-mesh peer, you may create your pull request with your new peer or use existing ones: https://github.com/vikulin/public-peers

## Known issues

### 1. Log message:
```
An error occurred starting multicast: listen udp6 [::]:9001: socket: address family not supported by protocol
```
and
```
An error occurred starting TUN/TAP: operation not supported
```

### Caused by:
The device has no IPv6 support


### 2. Log message:
```
An error occurred starting TUN/TAP: permission denied
```

### Caused by:
IPv6 support is not enabled. See the solution: https://github.com/yggdrasil-network/yggdrasil-go/issues/479#issuecomment-519512395

### 3. Mesh infinite output in log:
 Connected SCTP ...
 
 Disconnected SCTP ...

### Caused by:
Docker interface docker0 is conflicting with SCTP bind process. The issue can be resolved by removing docker.

## License

RiV-mesh is released under the **GNU Lesser General Public License v3.0 (LGPLv3)**:

- **License**: LGPLv3 (standard version, no special exceptions)
- **Commercial Use**: Allowed with full LGPLv3 compliance requirements
- **Source Code**: Available at [GitHub](https://github.com/vikulin/RiV-mesh)
- **Compliance**: Users must provide source code and comply with LGPLv3 terms

### Key LGPLv3 Requirements:
- **Source Code Sharing**: Must provide v6Space™ source code when distributing
- **Application Code Sharing**: Must share source code of applications using RiV-mesh
- **Build Instructions**: Must provide installation and build information
- **Modifications**: Any modifications to RiV-mesh must remain open source

**Note**: RiV-mesh is the legacy version. v6Space™ uses a dual-license model with additional UI components under CC BY-NC 4.0.
