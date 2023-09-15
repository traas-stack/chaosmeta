# docker build -t registry.cn-hangzhou.aliyuncs.com/chaosmeta/chaosmetad-demo:v0.2.0 -f chaosmetad-demo.Dockerfile .
From centos:centos7
ADD ./jdk-8u361-linux-x64.tar.gz /usr/local
RUN yum install -y iproute && yum clean all
ENV CHAOSMETAD_VERSION=0.2.0
ADD ./chaosmetad-$CHAOSMETAD_VERSION.tar.gz /opt/chaosmeta
RUN echo 'export JAVA_HOME=/usr/local/jdk1.8.0_361' >> /etc/profile && echo 'export PATH=$PATH:$JAVA_HOME/bin' >> /etc/profile && echo 'export PATH=$PATH:/opt/chaosmeta/chaosmetad-'${CHAOSMETAD_VERSION} >> /etc/profile
