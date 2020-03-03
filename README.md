# SAP Host Exporter

This is a bespoke Prometheus exporter enabling the monitoring of SAP systems (a.k.a. SAP NetWeaver applications).

[![Build Status](https://travis-ci.org/SUSE/sap_host_exporter.svg?branch=master)](https://travis-ci.org/SUSE/sap_host_exporter)

## Table of Contents
1. [Features](#features)
2. [Installation](#installation)
3. [Usage](#usage)
4. [Development](doc/devel.md)
5. [License](#license)

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

T.B.D.

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
