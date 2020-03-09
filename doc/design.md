# Design Notes

This document describes the rationale behind design decisions taken during the development of this project.

## Goals

- Export runtime statistics about the various SAP cluster components from existing data sources, to be consumed by a Prometheus monitoring stack.

## Non-goals

- Maintain an internal, consistent, persisting representation of the cluster state; since the original source of truth is distributed, we want to avoid the complexity of a stateful middleware.


## Structure

The project consist in a small HTTP application that exposes runtime data in a line protocol.
  
A series of "metric collectors" are consumed by the main application entry point, [`main.go`](../main.go), where they are registered with the Prometheus client and then exposed via its HTTP handler.

Concurrency is handled internally by a worker pool provided by the Prometheus library, but this implementation detail is completely obfuscated to the consumers.

The data sources are read every time an HTTP request comes, and the collected metrics are not shared: their lifecycle corresponds with the request's.

The `internal` package contains common code shared among all the other packages, but not intended for usage outside this projects.

## Collectors

Inside the `collector` package, you wil find the code of the main logic of the project: these are a number of [`prometheus.Collector`](https://github.com/prometheus/client_golang/blob/b25ce2693a6de99c3ea1a1471cd8f873301a452f/prometheus/collector.go#L16-L63) implementations, one for each cluster component (that we'll call _subsystems_), like the Start Service or the Enqueue Server.

Common functionality is provided by composing the [`DefaultCollector`](../collector/default_collector.go). 

Each subsystem collector has a dedicated package; some are very simple, some are little more nuanced. 

The collectors usually consume a SOAP web service called SAPControl and the retrieved data is then used to build the Prometheus metrics. 

More details about the metrics themselves can be found in the [metrics](metrics.md) document.
