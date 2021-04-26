# Developer notes

1. [Makefile](#makefile)
2. [Generated Code](#generated-code)
3. [OBS packaging](#obs-packaging)
4. [SAP learning material](#sap-learning-material)


## Makefile

Most development tasks can be accomplished via [make](../Makefile).

For starters, you can run the default target with just `make`.

The default target will clean, analyse, test and build the amd64 binary into the `build/bin` directory.

You can also cross-compile to the various architectures we support with `make build-all`.


## Generated code

Some of the code used in this repository is automatically generated.

You can use the `make generate` target (which in turn runs `go generate`) to update generated code by taking advantage of the `//go:generate` annotation in the code.

### Mocks

We generate the mocks with the [GoMock](https://github.com/golang/mock) library. 

All the mocked packages should follow the same convention and be put in the corresponding `mock_*` package inside the `test` directory.

Only public interfaces should need to be mocked.

### SAPControl web service

The for the [SAPControl web service](lib/sapcontrol/soap_wsdl.go), we generated the basic structure with [hooklift/gowsdl](https://github.com/hooklift/gowsdl), then extracted and adapted only the parts of the web service that we actually need.

For reference, you can find the full, generated, web service code [here](_generated_soap_wsdl.go), but bear in mind that we don't intend to use its generated code as it is. As such, note that this file is not covered by the `make generate` target.


## OBS Packaging

The CI will automatically publish GitHub releases to SUSE's Open Build Service: to perform a new release, just publish a new GH release. Tags must always follow the [SemVer](https://semver.org/) scheme.

If you wish to produce an OBS working directory locally, having configured [`osc`](https://en.opensuse.org/openSUSE:OSC) already, you can run:
```
make obs-workdir
```
This will checkout the OBS project and prepare a new OBS commit in the `build/obs` directory.

You can use the `OSB_PROJECT`, `REPOSITORY`, `VERSION` and `REVISION` environment variables to change the behaviour of OBS-related make targets.

By default, the current Git working directory is used to infer the values of `VERSION` and `REVISION`, which are used by OBS source services to generate a compressed archive of the sources.  

For example, if you were on a feature branch of your own fork, you may want to change these variables, so:
```bash
git checkout feature/xyz
git push johndoe feature/xyz # don't forget to push changes your own fork remote
export OBS_PROJECT=home:JohnDoe
export REPOSITORY=johndoe/sap_host_exporter
make exporter-obs-workdir
``` 
will prepare to commit on OBS into `home:JohnDoe/sap_host_exporter` by checking out the branch `feature/yxz` from `github.com/johndoe/sap_host_exporter`.

At last, to actually perform the commit into OBS, run: 
```bash
make exporter-obs-commit
``` 

Note that that actual continuously deployed releases also involve an intermediate step that updates the changelog automatically with the markdown text of the GitHub release.


## SAP learning material

This section will provide some initial pointers to understand the documentation material behind SAP NetWeaver.

Since we don't control the sources, some small changes may be introduced in the future, and the links might stop working; please feel free to submit a PR in case the documentation becomes outdated.

### Exploring the SAPControl web service

You can find the full SAPControl official documentation [here](https://www.sap.com/documents/2016/09/0a40e60d-8b7c-0010-82c7-eda71af511fa.html).

In order to learn the SOAP interface, you can use the following Python script (an example and adapted extracted from the previous link):

```python
#!/usr/bin/python3

from suds.client import Client
# Create proxy from WSDL
SAP_URL = 'http://host:port?wsdl'
client = Client(SAP_URL)
# Call unprotected webmethod with complex output
print("Get process list \n")

result = client.service.GetProcessList()
print(result)
print("PID of First process: \n")
# Access output data
print('PID:', result[0][0].pid)
```
