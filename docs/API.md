# Metricator's API

This page describes Metricator's API.

## URL forming

API URLs are formed like: `/api/vX/TYPE/OTHER_DATA`, where:

* `vX` is a version (e.g. `v1`).
* `TYPE` - request type (see below).
* `OTHER_DATA` - other data that is needed to properly process request of `TYPE`.

## Versions

Metricator's API currently supports API version 1.

## Handlers

### Get Metricator's build info

**URL**: `/api/vX/info`

Returns a JSON with build info with following items:

* `Branch` for git branch.
* `Build` for commit **number**.
* `CommitHash` for git commit hash.
* `Version` for tagged (or not tagged) build.

Caveats:

1. `Version` contains string formed with [Semantic Versioning](https://semver.org/spec/v2.0.0.html) in mind, so there is either "X.Y.Z" for tagged release or "X.Y.Z-dev` for development (built not from tag) release.

Example:

```json
{"Branch":"master","Build":"24","CommitHash":"095e24a2b0b908d790153a4ed0bee51d9dba7685","Version":"0.1.0-dev"}
```

### Get registered applications list

**URL**: `/api/vX/apps_list`

Returns a JSON with all known applications.

Example:

```json
["app1", "app2", "gw1", "gw2"]
```

### Get all known metrics for application

**URL**: `/api/vX/metrics/APP_NAME`

Where:

* `APP_NAME` is a name of registered application.

Example:

```json
[
    {
        "Name":"dnsdist_frontend_tcpgaveup/frontend:127.0.0.1:53/proto:TCP/thread:0",
        "Description":"Amount of TCP connections terminated after too many attempts to get a connection to the backend",
        "Type":"counter",
        "Value":"0",
        "Params":null
    },
    {
        "Name":"dnsdist_frontend_queries/frontend:127.0.0.1:53/proto:UDP/thread:0",
        "Description":"Amount of queries received by this frontend",
        "Type":"counter",
        "Value":"0",
        "Params":null
    },
    {
        "Name":"dnsdist_frontend_queries/frontend:127.0.0.1:53/proto:UDP/thread:0",
        "Description":"Amount of queries received by this frontend",
        "Type":"counter",
        "Value":"53842",
        "Params":null,
    }
]
```

### Get single metric data

**URL**: `/api/vX/metrics/APP_NAME/METRIC_NAME`

Where:

* `APP_NAME`  is a name of registered application.
* `METRIC_NAME` is a name of metric exposed by this application.
