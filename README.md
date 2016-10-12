lmsensors_exporter [![Build Status](https://travis-ci.org/mdlayher/lmsensors_exporter.svg?branch=master)](https://travis-ci.org/mdlayher/lmsensors_exporter) [![GoDoc](http://godoc.org/github.com/mdlayher/lmsensors_exporter?status.svg)](http://godoc.org/github.com/mdlayher/lmsensors_exporter)
==================

Command `lmsensors_exporter` provides a Prometheus exporter for
[lm_sensors](https://en.wikipedia.org/wiki/Lm_sensors) sensor metrics.
MIT Licensed.

Usage
-----

Available flags for `lmsensors_exporter` include:

```
$ ./lmsensors_exporter -h
Usage of ./lmsensors_exporter:
  -telemetry.addr string
        address for lmsensors exporter (default ":9165")
  -telemetry.path string
        URL path for surfacing collected metrics (default "/metrics")
```
