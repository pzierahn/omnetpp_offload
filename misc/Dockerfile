FROM ubuntu:latest

#
# Setup
#

ENV WORKDIR=/install
WORKDIR $WORKDIR
ENV DEBIAN_FRONTEND noninteractive
ENV MODE release

#
# Get essential dependencies
#

RUN ln -fs /usr/share/zoneinfo/Europe/Berlin /etc/localtime

RUN apt-get update; \
    apt-get upgrade -y; \
    apt-get dist-upgrade -y; \
    apt-get install -y --no-install-recommends git nano wget tzdata apt-utils; \
    apt-get install -y --no-install-recommends build-essential clang make bison flex libxml2-dev zlib1g-dev python; \
    apt-get clean; \
    rm -rf /var/lib/apt/lists/*

#
# Get build omnetpp
#

ARG VERSION=5.6.2
# ARG VERSION=6.0pre11

RUN wget https://github.com/omnetpp/omnetpp/releases/download/omnetpp-$VERSION/omnetpp-$VERSION-src-core.tgz --no-check-certificate -O omnetpp-src-core.tgz; \
    tar xf omnetpp-src-core.tgz; \
    rm omnetpp-src-core.tgz; \
    mv omnetpp-$VERSION omnetpp

ENV PATH $WORKDIR/omnetpp/bin:$PATH

RUN cd omnetpp; \
    ./configure WITH_QTENV=no WITH_OSG=no WITH_OSGEARTH=no; \
    make -j $(nproc) MODE=release base
