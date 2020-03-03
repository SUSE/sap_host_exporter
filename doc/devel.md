# Developing

# Learning materials:

This guide will provide some guided to understand the material behind Netweaver.
The link might be broken or some small changes will needed during the time, just submit a PR to fix if you find broken links.

This guide goal is to provide some entrypoints.


##  Exploratory/Learning SOAP of SapControl interface:

You can find the full documentation here: https://www.sap.com/documents/2016/09/0a40e60d-8b7c-0010-82c7-eda71af511fa.html?infl=547aee22-5559-4c30-9b86-3e2df9e5d806



In order to learn the interface of SOAP you can use the folowing sapcontrol script: ( example extracted from previous link and adapted)

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
