# Metricator Client

Metricator client was created to help with Metricator daemon communication. It is able to produce different output based on selected format.

## Parameters

| Parameter | Type | Description |
| --------- | ---- | ----------- |
| `-application` | string | Name of application to query. |
| `-apps-list` | bool | Request type: applications list. |
| `-config` | string | Path to configuration file. **MANDATORY** |
| `-metric` | string | Request type: single metric. Name of metric to request. |
| `-metricator-host` | URL | URL to Metricator daemon. (e.g. `http://127.0.0.1:34421`). **MANDATORY** |
| `-metricator-timeout` | integer | Timeout in seconds for Metricator Client's HTTP requests. By default - 5 seconds. |
| `-metrics-list` | bool | Request type: list of metrics. **Requires `-application` parameter to be filled.** |
| `-output` | string | Type of output to produce (see below). |

## Request types

One of following parameters should be defined:

* `-apps-list` to get listing of applications that is registered at Metricator daemon.
* `-metrics-list` with `-application` parameters to get list of metrics for application which Metricator know.
* `-metric` with name of metric and `-application` (with a name of application) parameters to get specific metric for application which Metricator know.

See Examples section below.

## Outputs

Currently Metricator client is able to produce JSON and "Plain By Line" outputs. Their meanings:

* When `-output=json` is specified (or `-output` wasn't specified at all) Metricator Client will just dump response from Metricator Daemon.
* When `-output=plain-by-line` is specified Metricator Client will transform received data into line-by-line output, e.g. every application name on separate line, every metric name on separate line, etc.

## Examples

* Get list of applications registered at Metricator daemon line-by-line (for later use with metrics autodiscovery helper):

```shell
metricator-client -config=./metricator.yaml -metricator-host http://127.0.0.1:34421 -output plain-by-line -apps-list
```

* Get list of metrics for application `test` in JSON format:

```shell
metricator-client -config ./metricator.yaml -metricator-host http://127.0.0.1:34421 -application test -metrics-list
```

* Get specific metric for application `test`:

```shell
metricator-client -config ./metricator.yaml -metricator-host http://127.0.0.1:34421 -application test -metric mymegametric
```
