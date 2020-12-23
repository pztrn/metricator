# Installing Metricator

There are three ways to install metricator:

1. Using Docker image.
2. Using pre-built binaries.
3. From source

This page describes all of them.

## Docker

To run Metricator in Docker simply run:

```bash
docker run -v $(pwd)/metricator.yaml:/config.yaml -p 8080:34421 registry.gitlab.pztrn.name/pztrn/metricator:latest
```

Don't forget to create configuration file as described [here](CONFIGURE.md)!

### docker-compose

```docker-compose
version: "2.4"

services:
  metricator:
    image: registry.gitlab.pztrn.name/pztrn/metricator:latest
    restart: always
    ports:
      - 8080:34421
    volumes:
      - /full/path/to/metricator.yaml:/config.yaml
```

## Pre-built binaries

To be written. Stay tuned.

## From source

To be written. Stay tuned.
