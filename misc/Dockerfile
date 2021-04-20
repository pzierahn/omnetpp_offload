FROM ubuntu:latest

#
# Setup
#

ENV WORKDIR=/install
WORKDIR $WORKDIR
ENV DEBIAN_FRONTEND noninteractive

#
# Get essential dependencies
#

RUN apt-get update; \
    apt-get upgrade -y; \
    apt-get dist-upgrade -y; \
    apt-get install -y --no-install-recommends git nano wget tzdata; \
    apt-get install -y --no-install-recommends build-essential clang make bison flex git libxml2-dev zlib1g-dev; \
    apt-get clean; \
    rm -rf /var/lib/apt/lists/*

RUN ln -fs /usr/share/zoneinfo/Europe/Berlin /etc/localtime; \
    dpkg-reconfigure --frontend noninteractive tzdata

# RUN apt-get install -y --no-install-recommends build-essential clang make bison flex git libxml2-dev zlib1g-dev

#
# Get build omnetpp
#

ARG VERSION=5.6.2
# ARG VERSION=6.0pre10

RUN wget https://github.com/omnetpp/omnetpp/releases/download/omnetpp-$VERSION/omnetpp-$VERSION-src-core.tgz --no-check-certificate -O omnetpp-src-core.tgz; \
    tar xf omnetpp-src-core.tgz; \
    rm omnetpp-src-core.tgz; \
    mv omnetpp-$VERSION omnetpp

ENV PATH $WORKDIR/omnetpp/bin:$PATH

RUN cd omnetpp; \
    ./configure WITH_QTENV=no WITH_OSG=no WITH_OSGEARTH=no; \
    make -j $(nproc) MODE=release base

# ENV PATH /root/omnetpp/bin:$PATH
