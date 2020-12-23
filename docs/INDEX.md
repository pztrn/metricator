# Metricator

Welcome to Metricator docs.

## What's the idea

Main idea for Metricator is to provide high-speed proxy between classic NMSes (like NetXMS, Nagios and so on) and other software that exposes metrics in Prometheus format.

Why proxy? Let's go by example.

Imagine that software you wish to monitor with classic NMS exposes 250 metric items and you wrote a simple script which returns needed data. Classic NMSes is able to process only one value per checker (usually), which means that processing every metric will do 250 HTTP requests to monitored software which is obviously not good.

Of course you can write own script that will request once and another script that will parse saved data, but how would you deal with parametrized metrics (like `dnsdist_frontend_responses{frontend="127.0.0.1:53",proto="UDP",thread="0"}`)? This will definetely took a lot of brain cells and time.

Metricator instead acts like Prometheus itself from monitored software PoV: it performs one HTTP request each timeout (configured in config), parses data and making it available to other requesters. As parsed metric data stored in memory it will be blazing fast and won't overload your system.

Also Metricator "reformats" metric names and parameters to be more easily parsed if needed, so:

```prometheus
dnsdist_frontend_responses{frontend="127.0.0.1:53",proto="UDP",thread="0"}
```

became:

```prometheus
dnsdist_frontend_responses/frontend:127.0.0.1:53/proto:UDP/thread:0
```

## Docs

* [Installation](INSTALL.md)
* [Configuration](CONFIGURE.md)
* [API](API.md)
