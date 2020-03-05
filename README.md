# SAP Host Exporter

This is a bespoke Prometheus exporter enabling the monitoring of SAP systems (a.k.a. SAP NetWeaver applications).

[![Build Status](https://travis-ci.org/SUSE/sap_host_exporter.svg?branch=master)](https://travis-ci.org/SUSE/sap_host_exporter)

## Table of Contents
1. [Features](#features)
2. [Installation](#installation)
3. [Usage](#usage)
4. [Development](doc/devel.md)
5. [Contributing](#contributing)
6. [License](#license)

## Features

T.B.D.

## Installation

The project can be installed in many ways, including but not limited to:

1. [Manual clone & build](#manual-clone-&-build)
2. [Go](#go)
3. [RPM](#rpm)


### Manual clone & build

```
git clone https://github.com/SUSE/sap_host_exporter
cd sap_host_exporter
make
make install
```

### Go

```
go get github.com/SUSE/sap_host_exporter
```

### RPM
You can find the repositories for RPM based distributions in [SUSE's Open Build Service](https://build.opensuse.org/package/show/server:monitoring/prometheus-sap_host_exporter).  
On openSUSE or SUSE Linux Enterprise you can just use the `zypper` system package manager:
```shell
export DISTRO=SLE_15_SP2 # change as desired
zypper addrepo https://download.opensuse.org/repositories/server:/monitoring/$DISTRO/server:monitoring.repo
zypper install prometheus-sap_host_exporter
```

## Usage

You can run the exporter in any of the netweaver cluster nodes.

```
$ ./sap_host_exporter --
INFO[0000] Serving metrics on 0.0.0.0:9680
```

Though not strictly required, it is _strongly_ advised to run it in all the nodes.

It will export the metrics under the `/metrics` path, on port `9680` by default.

While the exporter can run outside a Netweaver cluster node, it won't export any metric it can't collect.
A warning message will inform the user of such cases.

**Hint:**
You can deploy a full Netweaver Cluster via Terraform with [SUSE/ha-sap-terraform-deployments](https://github.com/SUSE/ha-sap-terraform-deployments), and also monitoring it.

### Configuration

All the runtime parameters can be configured either via CLI flags or via a configuration file, both or which are completely optional.

For more details, refer to the help message via `sap_host_exporter --help`.

**Note**:
the built-in defaults are tailored for the latest version of SUSE Linux Enterprise and openSUSE.

The program will scan, in order, the current working directory, `$HOME/.config`, `/etc` and `/usr/etc` for files named `sap_host_exporter.(yaml|json|toml)`.
The first match has precedence, and the CLI flags have precedence over the config file.

Please refer to the example YAML configuration for more details.

### systemd integration

A [systemd unit file](ha_cluster_exporter.service) is provided with the RPM packages. You can enable and start it as usual:

```
systemctl --now enable prometheus-sap_host_exporter
```


# Contribuiting

Pull requests are more than welcome!

We recommend having a look at the [design document](doc/design.md) before contributing.

Also for learning material take a look at [devel notes](doc/devel.md)


## License

Copyright 2019-2020 SUSE LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
