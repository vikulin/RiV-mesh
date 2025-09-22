# RiV-mesh - Decentralized IPv6 Mesh Network

## Introduction

RiV-mesh is the legacy implementation of a fully end-to-end encrypted IPv6
mesh network, designed to provide secure connectivity between a wide spectrum of endpoint devices like IoT devices,
desktop computers or even routers.
It is lightweight, self-arranging, supported on multiple
platforms and allows pretty much any IPv6-capable application
to communicate securely with other network nodes.

**Note**: RiV-mesh is now maintained for legacy compatibility. For new deployments, we recommend using [v6Space™](https://github.com/RiV-chain/v6Space) - the next-generation mesh networking platform with enhanced features, modern architecture, agentic capabilities, and MCP (Model Context Protocol) tools for AI agent communication and autonomous operations.

## Relationship to v6Space™

RiV-mesh is the **legacy version** of the mesh networking technology that has evolved into **v6Space™**:

### RiV-mesh (Legacy)
- **Status**: Legacy maintenance mode
- **License**: LGPLv3 only
- **Features**: Core mesh networking functionality
- **Use Case**: Existing deployments, compatibility

### v6Space™ (Current)
- **Status**: Active development
- **License**: Dual-license (LGPLv3 + CC BY-NC 4.0)
- **Features**: Enhanced mesh networking + modern UI + MCP support + agentic capabilities
- **Agentic Features**: AI agent communication, autonomous operations, Web3 integration
- **MCP Tools**: Model Context Protocol integration for AI model communication
- **Use Case**: New deployments, modern applications, AI agent networks, Web3 infrastructure

**Migration**: RiV-mesh configurations are compatible with v6Space™. Simply replace the binary and continue using your existing `mesh.conf` files.

## Core Functionality

RiV-mesh provides the decentralized IPv6 mesh overlay network infrastructure for applications like CupLink™, which requires a VPN service to establish the mesh network environment for peer-to-peer communication.

## Supported Platforms

RiV-mesh works on a number of platforms, including Linux, macOS, Ubiquiti
EdgeRouter, VyOS, Windows, FreeBSD, OpenBSD and OpenWrt.

Please see our [Installation](https://github.com/RiV-chain/RiV-mesh-builds#riv-mesh-build) 
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

Other OS packages can be built in this repo: https://github.com/RiV-chain/RiV-mesh-builds.

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

Documentation is available [on our website](https://riv-chain.github.io/RiV-mesh/).

- [Installing RiV-mesh](https://riv-chain.github.io/RiV-mesh/)
- [Configuring RiV-mesh](https://riv-chain.github.io/RiV-mesh/)
- [Frequently asked questions](https://riv-chain.github.io/RiV-mesh/)
- [Version changelog](CHANGELOG.md)

## Work in progress:

<img width="253" alt="DDNS" src="https://github.com/RiV-chain/RiV-mesh/assets/743622/05757b9f-2053-4503-9037-24e5ff992554">

## Company Information

**RiV-mesh** is developed by **RiV Chain™ Limited**, an Ireland-based FinTech and infrastructure R&D company focused on decentralized networking and agentic systems.

- **Company**: RiV Chain™ Limited
- **Location**: Ireland
- **Focus**: Decentralized mesh networking, AI agent coordination, distributed identity, Web3 infrastructure
- **Specialization**: Agentic Web technologies, MCP tools, autonomous AI agent communication
- **Website**: [rivchain.org](https://rivchain.org)
- **Contact**: For technical support and licensing inquiries



## Public peers
If you are operating a RiV-mesh peer, you may create your pull request with your new peer or use existing ones: https://github.com/RiV-chain/public-peers

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
- **Source Code**: Available at [GitHub](https://github.com/RiV-chain/v6Space)
- **Compliance**: Users must provide source code and comply with LGPLv3 terms

### Key LGPLv3 Requirements:
- **Source Code Sharing**: Must provide v6Space™ source code when distributing
- **Application Code Sharing**: Must share source code of applications using RiV-mesh
- **Build Instructions**: Must provide installation and build information
- **Modifications**: Any modifications to RiV-mesh must remain open source

For complete licensing details, see:
- [Full LGPLv3 License](LICENSE)
- [v6Space™ License Information](https://rivchain.org/legal/license.html)
- [Component Separation Guide](https://rivchain.org/legal/component-separation.html)

**Note**: RiV-mesh is the legacy version. v6Space™ uses a dual-license model with additional UI components under CC BY-NC 4.0.
