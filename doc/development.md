# Developer notes

## Generated code

Some of the code used in this repository is automatically generated.

### SAPControl web service

The for the [SAPControl web service](internal/sapcontrol/soap_wsdl.go), we generated the basic structure with [hooklift/gowsdl](https://github.com/hooklift/gowsdl), then extracted and adapted only the parts of the web service that we actually need.

For reference, you can find the full, generated, web service code [here](_generated_soap_wsdl.go), but bear in mind that we don't intend to use its generated code as it is.

### Mocks

We generate the mocks with the [GoMock](https://github.com/golang/mock) library.

For example, the [`sapcontrol.WebService`](internal/sapcontrol/webservice.go) interface is mocked with the following command:
```
mockgen --source ./internal/sapcontrol/webservice.go WebService > test/mock_sapcontrol/webservice.go
```

All the mocked packages should follow the same convention and be put in the corresponding `mock_*` package inside the `test` directory.

Only public interfaces should need to be mocked.

## SAP learning material

This section will provide some initial pointers to understand the documentation material behind SAP NetWeaver.

Since we don't control the sources, some small changes may be introduced in the future, and the links might stop working; please feel free to submit a PR in case the documentation becomes outdated.

###  Exploring the SAPControl web service

You can find the full documentation here: https://www.sap.com/documents/2016/09/0a40e60d-8b7c-0010-82c7-eda71af511fa.html

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
