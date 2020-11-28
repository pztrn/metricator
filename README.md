# Metricator

Simple proxy between prometheus-powered application and your NMS.

## Why

I'm using NetXMS to monitor all of my systems. Many things I use exports metrics in prometheus format which can be utilized by custom parsing script. But I've encounter a performance problem when I need to monitor 30 metrics - parsing script will make 30 requests to prometheus endpoint which might affect performance.

Metricator will issue only one request and cache data in memory between them. Also it will expose HTTP API to get single metric which can be easily utilized with any NMS.

## Caveats

* No authorization. **DO NOT** expose Metricator to wild world!

## Installation

TBW

## Configuration

TBW

## Documentation

Check [docs directory](/docs/INDEX.md).
