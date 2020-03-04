# Developer notes

# Generatitng structs

go get github.com/hooklift/gowsdl@506189e0354e7197035e46515b80a0836944acf2



## Learning materials

This guide will provide some initial pointers to understand the material behind Netweaver.
The link might be broken or some small changes will be needed in the future; please feel free to submit a PR to fix any broken links you find.


##  Exploring and Learning the SOAP SapControl interface

You can find the full documentation here: https://www.sap.com/documents/2016/09/0a40e60d-8b7c-0010-82c7-eda71af511fa.html?infl=547aee22-5559-4c30-9b86-3e2df9e5d806



In order to learn the SOAP interface, you can use the following Python script: (example extracted from the previous link and adapted)

```python
#! /usr/bin/python3

from suds.client import Client
# Create proxy from WSDL
SAP_URL = 'http://10.162.32.183:50013/?wsdl'
client = Client(SAP_URL)
# Call unprotected webmethod with complex output
print("Get process list \n")

result = client.service.GetProcessList()
print(result)
print("PID of First process: \n")
# Access output data
print('PID:', result[0][0].pid)
```
