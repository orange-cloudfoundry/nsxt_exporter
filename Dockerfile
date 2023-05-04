FROM        ubuntu:latest
MAINTAINER  Xavier MARCELET <xavier.marcelet@orange.com>

ADD nsxt_exporter /bin/nsxt_exporter

ENTRYPOINT ["/bin/nsxt_exporter"]
EXPOSE     8080
