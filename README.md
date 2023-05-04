# Introduction

This project implements an [prometheus] exporter for vmware [NSX-T]. It provides
metrics for `Cluster`, `LBService`, `LBPool`, `LBVirtualServer`, `Tier0` and `Tier1` objects.

This exporter is not suitable for NSX-T admins that want to achieve a full monitoring their
infrastructure but can be handy for NSX-T users that need to implement a lightweight monitoring
of their hosting environment.


[prometheus]: https://prometheus.io/docs/instrumenting/exporters/
[NSX-T]: https://developer.vmware.com/apis/1248/nsx-t

# NSX Version

This current implementation has been designed and tested for NSX-T API `3.2.1`.

# Metrics

The document [METRICS.md](./METRICS.md) describe all generated metrics.


# Install

* from release assets
  * download assets for your architecture from [latest release]
  * extract tarball: `tar xzf nsxt_exporter_$version_linux_amd64.tar.gz`
* from go install: `go install gihub.com/orange-cloudfoundry/nsxt_exporter`
* from source: `CGO_ENABLED=0 go build -o nsxt_exporter -ldflags='-s -w' .`

[latest release]: https://github.com/orange-cloudfoundry/nsxt_exporter/releases

# Usage

```
usage: nsxt_exporter [<flags>]


Flags:
  -h, --[no-]help          Show context-sensitive help (also try --help-long and --help-man).
      --config=config.yml  Configuration file path
      --[no-]version       Show application version.
```

# Configuration

See [config.yml.sample](./config.yml.sample)

# Development

- static check analysis:
  - install: https://golangci-lint.run/usage/install/#local-installation
  - run: `golangci-lint run --config .golangci.yml`
