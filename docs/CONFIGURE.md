# Configuration

This page describes every configuration value.

## The Table

| Variable | Type | Description |
| -------- | ---- | ----------- |
| applications > APPNAME > endpoint | string | Prometheus metrics endpoint URL. |
| applications > APPNAME > headers > MAP | map | Headers which should be added to request. See example below. |
| applications > APPNAME > time_between_requests | string | time.Duration-compatible string which represents timeout between requests. |

Where:

* `APPNAME` is a custom name for application.

## Headers map example

```yaml
applications:
  example:
    headers:
      HeaderOne: valueOne
      HeaderTwo: valueTwo
```
