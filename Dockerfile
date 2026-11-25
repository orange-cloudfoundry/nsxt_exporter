FROM        ubuntu:latest
MAINTAINER  Xavier MARCELET <xavier.marcelet@orange.com>

ARG TARGETPLATFORM
COPY $TARGETPLATFORM/nsxt_exporter /bin/nsxt_exporter

ENTRYPOINT ["/bin/nsxt_exporter"]
EXPOSE     8080
