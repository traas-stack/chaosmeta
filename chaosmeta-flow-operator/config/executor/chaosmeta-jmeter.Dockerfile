FROM centos:centos7
WORKDIR /opt
ADD jdk-8u361-linux-x64.tar.gz /usr/local/
ADD apache-jmeter-5.6.tgz /opt/
ENV JAVA_HOME=/usr/local/jdk1.8.0_361 \
    PATH=/usr/local/jdk1.8.0_361/bin:/opt/apache-jmeter-5.6/bin:$PATH
