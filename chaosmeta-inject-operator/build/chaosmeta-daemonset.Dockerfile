# docker build -t ghcr.io/traas-stack/chaosmeta-daemon:v0.1.0 -f chaosmeta-daemonset.Dockerfile .
From centos:centos7
ENV CHAOSMETAD_VERSION=0.1.0
ADD ./chaosmetad-$CHAOSMETAD_VERSION.tar.gz /opt/chaosmeta
CMD while true; do if [ ! -d "/tmp/chaosmetad-$CHAOSMETAD_VERSION" ]; then cp -r /opt/chaosmeta/chaosmetad-$CHAOSMETAD_VERSION /tmp/chaosmetad-$CHAOSMETAD_VERSION; fi; sleep 600; done
