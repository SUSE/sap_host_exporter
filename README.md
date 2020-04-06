# SAP Host Exporter

This is a bespoke Prometheus exporter enabling the monitoring of SAP systems (a.k.a. SAP NetWeaver applications).

[![Build Status](https://travis-ci.org/SUSE/sap_host_exporter.svg?branch=master)](https://travis-ci.org/SUSE/sap_host_exporter)


## Table of Contents
1. [Features](#features)
2. [Installation](#installation)
3. [Usage](#usage)
   1. [Metrics](doc/metrics.md)
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

You can run the exporter as follows:

```shell
./sap_host_exporter --sap-control-url http://$SAP_HOST:$SAP_PORT
```

It will export the metrics under the `/metrics` path, on port `9680` by default.

Though not strictly required, it is advised to run it in the nodes of the cluster and access the SAPControl web service locally.

The exporter won't export any metric it can't collect, but since it doesn't care about which subsystems are present in the monitored target, failing to collect metrics is _not_ considered a hard failure condition.
Instead, in case some of the collectors fail to either register or perform collect cycles, a soft warning will be printed out in the log.

Refer to [doc/metrics.md](doc/metrics.md) for extensive details about all the exported metrics.

**Hint:**
You can deploy a full SAP NetWeaver cluster via Terraform with [SUSE/ha-sap-terraform-deployments](https://github.com/SUSE/ha-sap-terraform-deployments); 
this exporter and the whole Prometheus monitoring stack will be automatically installed and configured for you.

### Configuration

The runtime parameters can be configured either via CLI flags or via a configuration file, both or which are completely optional.

For more details, refer to the help message via `sap_host_exporter --help`.

**Note**:
the built-in defaults are tailored for the latest version of SUSE Linux Enterprise and openSUSE.

The program will scan, in order, the current working directory, `$HOME/.config`, `/etc` and `/usr/etc` for files named `sap_host_exporter.(yaml|json|toml)`.
The first match has precedence, and the CLI flags have precedence over the config file.

Please refer to the [example YAML configuration](doc/sap_host_exporter.yaml) for more details.

### systemd integration

A [systemd unit file](packaging/obs/prometheus-sap_host_exporter.spec) is provided with the RPM packages. You can enable and start it as usual:

```
systemctl --now enable prometheus-sap_host_exporter
```


## Contributing

Pull requests are more than welcome!

We recommend having a look at the [design document](doc/design.md) and the [development notes](doc/development.md) before contributing.


## License

Copyright 2020 SUSE LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
