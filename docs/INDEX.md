# Metricator

Welcome to Metricator docs.

## What's the idea

Main idea for Metricator is to provide high-speed proxy between classic NMSes (like NetXMS, Nagios and so on) and other software that exposes metrics in Prometheus format.

Why proxy? Let's go by example.

Imagine that software you wish to monitor with classic NMS exposes 250 metric items and you wrote a simple script which returns needed data. Classic NMSes is able to process only one value per checker (usually), which means that processing every metric will do 250 HTTP requests to monitored software which is obviously not good.

Metricator instead acts like Prometheus itself from monitored software PoV: it performs one HTTP request each timeout (configured in config), parses data and making it available to other requesters. As parsed metric data stored in memory it will be blazing fast and won't overload your system.
