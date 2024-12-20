# SAP Host Exporter

This is a bespoke Prometheus exporter enabling the monitoring of SAP systems (a.k.a. SAP NetWeaver applications).

[![Exporter CI](https://github.com/SUSE/sap_host_exporter/workflows/Exporter%20CI/badge.svg)](https://github.com/SUSE/sap_host_exporter/actions?query=workflow%3A%22Exporter+CI%22)

## Table of Contents

1. [Features](#features)
2. [Installation](#installation)
3. [Usage](#usage)
   1. [Configuration](#configuration)
   2. [Metrics](#metrics)
   3. [systemd integration](#systemd-integration)
5. [Contributing](#contributing)
   1. [Design](doc/design.md)
   2. [Development](doc/development.md)
6. [License](#license)

## Features

The exporter is a stateless HTTP endpoint. On each HTTP request, it pulls runtime data from the SAP system via the SAPControl web interface.

Exported data include:

- Start Service processes
- Enqueue Server stats
- AS Dispatcher work process queue stats  

## Installation

The project can be installed in many ways, including but not limited to:

1. [Manual clone and build](#manual-clone-and-build)
2. [Go](#go)
3. [RPM](#rpm)

### Manual clone and build

```shell
git clone https://github.com/SUSE/sap_host_exporter
cd sap_host_exporter
make build
make install
```

### Go

```shell
go get github.com/SUSE/sap_host_exporter
```

### RPM

You can find the repositories for RPM based distributions in [SUSE's Open Build Service](https://build.opensuse.org/package/show/devel:sap:monitoring:stable/prometheus-sap_host_exporter).  
On openSUSE or SUSE Linux Enterprise you can just use the `zypper` system package manager:

```shell
zypper install prometheus-sap_host_exporter
```

## Usage

You can run the exporter as follows:

```shell
./sap_host_exporter --sap-control-url $SAP_HOST:$SAP_CONTROL_PORT
```

Though not strictly required, it is advised to run the exporter locally in the target SAP instance host, and connect to the SAPControl web service via Unix Domain Sockets:

```shell
./sap_host_exporter --sap-control-uds /tmp/.sapstream50013
```

For further details on SAPControl, please refer to the [official SAP docs](https://www.sap.com/documents/2016/09/0a40e60d-8b7c-0010-82c7-eda71af511fa.html) to properly connect to the SAPControl service.

The exporter will expose the metrics under the `/metrics` path, on port `9680` by default.

### Configuration

The runtime parameters can be configured either via CLI flags or via a configuration file, both of which are completely optional.

For more details, refer to the help message via `sap_host_exporter --help`.

**Note**:
the built-in defaults are tailored for the latest version of SUSE Linux Enterprise and openSUSE.

The program will scan, in order, the current working directory, `$HOME/.config`, `/etc` and `/usr/etc` for files named `sap_host_exporter.(yaml|json|toml)`.
The first match has precedence, and the CLI flags have precedence over the config file.

Please refer to the [example YAML configuration](doc/sap_host_exporter.yaml) for more details.

### Metrics

The exporter won't export any metric it can't collect, but since it doesn't care about which subsystems are present in the monitored target, failing to collect metrics is _not_ considered a hard failure condition.
Instead, in case some of the collectors fail to either register or perform collect cycles, a soft warning will be printed out in the log.

Refer to [doc/metrics.md](doc/metrics.md) for extensive details about all the exported metrics.

### systemd integration

A [systemd unit file](packaging/obs/prometheus-sap_host_exporter.spec) is provided with the RPM packages. You can enable and start it as usual:

```
systemctl --now enable prometheus-sap_host_exporter
```

## Contributing

Pull requests are more than welcome!

We recommend having a look at the [design document](doc/design.md) and the [development notes](doc/development.md) before contributing.

## Copying

Copyright 2020-2025 SUSE LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this code repository except in compliance with the License.
You may obtain a copy of the License at

   <https://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
