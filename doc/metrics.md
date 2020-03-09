# Metrics

This document describes the metrics exposed by `sap_host_exporter`.

General notes:
- All the metrics are _namespaced_ with the prefix `sap`, which is followed by a _subsystem_, and both are in turn composed into a _Fully Qualified Name_ (FQN) of each metrics.
- All the metrics and labels _names_ are in snake_case, as conventional with Prometheus. That said, as much as we'll try to keep this consistent throughout the project, the label _values_ may not actually follow this convention, though (e.g. label is a hostname).

These are the currently implemented subsystems.

1. [SAP Start Service](#sap_start_service)
2. [SAP Enqueue Server](#sap_enqueue_server)


## SAP Start Service

The Start Service subsystem collects generic host-related metrics.

1. [`sap_start_service_processes`](#sap_start_service_processes) 

### `sap_start_service_processes`

#### Description

The processes started by the SAP Start Service.
For each process there will be a line with value `1`.

#### Labels

- `dispstatus`: the display status; one of `GRAY|GREEN|YELLOW|RED`.
- `name`: the name of the process.
- `pid`: the PID of the process.
- `textstatus`: a textual status of the process, e.g. `Running`

The total number of lines for this metric will be the cardinality of `pid`.

#### Example

```
# HELP sap_start_service_processes The processes started by the SAP Start Service
# TYPE sap_start_service_processes gauge
sap_start_service_processes{dispstatus="GREEN",name="enserver",pid="30787",textstatus="Running"} 1
sap_start_service_processes{dispstatus="GREEN",name="msg_server",pid="30786",textstatus="Running"} 1
```


## SAP Enqueue Server

The Enqueue Server (also known as the lock server) is the SAP system component that manages the lock table.

01. [`sap_enqueue_server_arguments_high`](#sap_enqueue_server_arguments_high)
02. [`sap_enqueue_server_arguments_max`](#sap_enqueue_server_arguments_max)
03. [`sap_enqueue_server_arguments_now`](#sap_enqueue_server_arguments_now)
04. [`sap_enqueue_server_arguments_state`](#sap_enqueue_server_arguments_state)
05. [`sap_enqueue_server_backup_requests`](#sap_enqueue_server_backup_requests)
06. [`sap_enqueue_server_cleanup_requests`](#sap_enqueue_server_cleanup_requests)
07. [`sap_enqueue_server_compress_requests`](#sap_enqueue_server_compress_requests)
08. [`sap_enqueue_server_dequeue_all_requests`](#sap_enqueue_server_dequeue_all_requests)
09. [`sap_enqueue_server_dequeue_errors`](#sap_enqueue_server_dequeue_errors)
10. [`sap_enqueue_server_dequeue_requests`](#sap_enqueue_server_dequeue_requests)
11. [`sap_enqueue_server_enqueue_errors`](#sap_enqueue_server_enqueue_errors)
12. [`sap_enqueue_server_enqueue_rejects`](#sap_enqueue_server_enqueue_rejects)
13. [`sap_enqueue_server_enqueue_requests`](#sap_enqueue_server_enqueue_requests)
14. [`sap_enqueue_server_lock_time`](#sap_enqueue_server_lock_time)
15. [`sap_enqueue_server_lock_wait_time`](#sap_enqueue_server_lock_wait_time)
16. [`sap_enqueue_server_locks_high`](#sap_enqueue_server_locks_high)
17. [`sap_enqueue_server_locks_max`](#sap_enqueue_server_locks_max)
18. [`sap_enqueue_server_locks_now`](#sap_enqueue_server_locks_now)
19. [`sap_enqueue_server_locks_state`](#sap_enqueue_server_locks_state)
20. [`sap_enqueue_server_owner_high`](#sap_enqueue_server_owner_high)
21. [`sap_enqueue_server_owner_max`](#sap_enqueue_server_owner_max)
22. [`sap_enqueue_server_owner_now`](#sap_enqueue_server_owner_now)
23. [`sap_enqueue_server_owner_state`](#sap_enqueue_server_owner_state)
24. [`sap_enqueue_server_replication_state`](#sap_enqueue_server_replication_state)
25. [`sap_enqueue_server_reporting_requests`](#sap_enqueue_server_reporting_requests)
26. [`sap_enqueue_server_server_time`](#sap_enqueue_server_server_time)
27. [`sap_enqueue_server_verify_requests`](#sap_enqueue_server_verify_requests)

### `sap_enqueue_server_arguments_high`

#### Example

```
# TYPE sap_enqueue_server_arguments_high gauge
sap_enqueue_server_arguments_high 104
```

### `sap_enqueue_server_arguments_max`

#### Example

```
# TYPE sap_enqueue_server_arguments_max gauge
sap_enqueue_server_arguments_max 56415
```

### `sap_enqueue_server_arguments_now`

#### Example

```
# TYPE sap_enqueue_server_arguments_now gauge
sap_enqueue_server_arguments_now 0
```

### `sap_enqueue_server_arguments_state`

#### Example

```
# TYPE sap_enqueue_server_arguments_state gauge
sap_enqueue_server_arguments_state 2
```

### `sap_enqueue_server_backup_requests`

#### Example

```
# TYPE sap_enqueue_server_backup_requests gauge
sap_enqueue_server_backup_requests 0
```

### `sap_enqueue_server_cleanup_requests`

#### Example

```
# TYPE sap_enqueue_server_cleanup_requests gauge
sap_enqueue_server_cleanup_requests 4
```

### `sap_enqueue_server_compress_requests`

#### Example

```
# TYPE sap_enqueue_server_compress_requests gauge
sap_enqueue_server_compress_requests 0
```

### `sap_enqueue_server_dequeue_all_requests`

#### Example

```
# TYPE sap_enqueue_server_dequeue_all_requests gauge
sap_enqueue_server_dequeue_all_requests 150372
```

### `sap_enqueue_server_dequeue_errors`

#### Example

```
# TYPE sap_enqueue_server_dequeue_errors gauge
sap_enqueue_server_dequeue_errors 0
```

### `sap_enqueue_server_dequeue_requests`

#### Example

```
# TYPE sap_enqueue_server_dequeue_requests gauge
sap_enqueue_server_dequeue_requests 85213
```

### `sap_enqueue_server_enqueue_errors`

#### Example

```
# TYPE sap_enqueue_server_enqueue_errors gauge
sap_enqueue_server_enqueue_errors 0
```

### `sap_enqueue_server_enqueue_rejects`

#### Example

```
# TYPE sap_enqueue_server_enqueue_rejects gauge
sap_enqueue_server_enqueue_rejects 4
```

### `sap_enqueue_server_enqueue_requests`

#### Example

```
# TYPE sap_enqueue_server_enqueue_requests gauge
sap_enqueue_server_enqueue_requests 109408
```

### `sap_enqueue_server_lock_time`

#### Example

```
# TYPE sap_enqueue_server_lock_time gauge
sap_enqueue_server_lock_time 174.574351
```

### `sap_enqueue_server_lock_wait_time`

#### Example

```
# TYPE sap_enqueue_server_lock_wait_time gauge
sap_enqueue_server_lock_wait_time 0
```

### `sap_enqueue_server_locks_high`

#### Example

```
# TYPE sap_enqueue_server_locks_high gauge
sap_enqueue_server_locks_high 104
```

### `sap_enqueue_server_locks_max`

#### Example

```
# TYPE sap_enqueue_server_locks_max gauge
sap_enqueue_server_locks_max 56415
```

### `sap_enqueue_server_locks_now`

#### Example

```
# TYPE sap_enqueue_server_locks_now gauge
sap_enqueue_server_locks_now 0
```

### `sap_enqueue_server_locks_state`

#### Example

```
# TYPE sap_enqueue_server_locks_state gauge
sap_enqueue_server_locks_state 2
```

### `sap_enqueue_server_owner_high`

#### Example

```
# TYPE sap_enqueue_server_owner_high gauge
sap_enqueue_server_owner_high 5
```

### `sap_enqueue_server_owner_max`

#### Example

```
# TYPE sap_enqueue_server_owner_max gauge
sap_enqueue_server_owner_max 56415
```

### `sap_enqueue_server_owner_now`

#### Example

```
# TYPE sap_enqueue_server_owner_now gauge
sap_enqueue_server_owner_now 0
```

### `sap_enqueue_server_owner_state`

#### Example

```
# TYPE sap_enqueue_server_owner_state gauge
sap_enqueue_server_owner_state 2
```

### `sap_enqueue_server_replication_state`

#### Example

```
# TYPE sap_enqueue_server_replication_state gauge
sap_enqueue_server_replication_state 2
```

### `sap_enqueue_server_reporting_requests`

#### Example

```
# TYPE sap_enqueue_server_reporting_requests gauge
sap_enqueue_server_reporting_requests 0
```

### `sap_enqueue_server_server_time`

#### Example

```
# TYPE sap_enqueue_server_server_time gauge
sap_enqueue_server_server_time 0
```

### `sap_enqueue_server_verify_requests`

#### Example

```
# TYPE sap_enqueue_server_verify_requests gauge
sap_enqueue_server_verify_requests 0
```
