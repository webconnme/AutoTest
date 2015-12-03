FROM ubuntu:10.04

RUN locale-gen ko_KR.UTF-8
RUN update-locale LANG=ko_KR.UTF-8
RUN dpkg-reconfigure locales

ENV LANG ko_KR.UTF-8
ENV LC_ALL ko_KR.UTF-8

ENV USER falinux

RUN \
    groupadd -g 1000 falinux && \
    useradd -m -p falinux -s /bin/bash -u 1000 -g falinux -G sudo falinux && \
    cd /home/falinux && \
    echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

RUN \
    sed -ri 's/\/archive\.ubuntu\.com/\/kr\.archive\.ubuntu\.com/g' /etc/apt/sources.list && \
    apt-get update && \
    apt-get -y install build-essential git-core pkg-config autoconf libtool unzip bc libncurses5-dev uuid-dev curl

###########################################
# Go 설치
###########################################
#https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz
WORKDIR /usr/local
RUN curl https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz | tar zx
ENV GOROOT /usr/local/go

ENV PATH $PATH:$GOROOT/bin/

ENV CGO_ENABLED 1

###########################################
# libsodium 0.7.1 설치
###########################################
WORKDIR /
RUN \
    curl https://download.libsodium.org/libsodium/releases/old/libsodium-0.7.1.tar.gz | tar zx && \
    cd /libsodium-0.7.1 && \
    ./configure --prefix=/usr && \
    make install && \
    ldconfig

###########################################
# zmq 4.1.2 준비
###########################################
# http://download.zeromq.org/zeromq-4.1.2.tar.gz
WORKDIR /
RUN \
    curl http://download.zeromq.org/zeromq-4.1.2.tar.gz | tar zx && \
    cd /zeromq-4.1.2 && \
    ./configure --prefix=/usr && \
    make install


###########################################
# go package 설치
###########################################
ENV GOPATH /vendor
RUN \
    mkdir /vendor && \
    cd /vendor && \
    go get -tags zmq_4_x github.com/alecthomas/gozmq && \
    go get github.com/satori/go.uuid && \
    go get github.com/go-martini/martini && \
    go get github.com/martini-contrib/binding && \
    go get github.com/martini-contrib/render && \
    go get github.com/martini-contrib/cors && \
    go get github.com/mikepb/go-serial

ENV GOPATH /vendor:/work
CMD ["/bin/bash"]

WORKDIR /work
