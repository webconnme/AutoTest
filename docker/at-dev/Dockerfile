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
ENV GOPATH=/work

ENV CGO_ENABLED 1

###########################################
# zmq 4.1.2 준비 ( 생략 )
###########################################
# http://download.zeromq.org/zeromq-4.1.2.tar.gz
WORKDIR /
RUN curl http://download.zeromq.org/zeromq-4.1.2.tar.gz | tar zx
WORKDIR /zeromq-4.1.2
RUN ./configure --prefix=/usr --with-libsodium=no && make install


###########################################
# go package 설치
###########################################
RUN \
    cd /usr/local/go && \
    go get -tags zmq_4_x github.com/alecthomas/gozmq && \
    go get github.com/satori/go.uuid && \
    go get github.com/go-martini/martini && \
    go get github.com/martini-contrib/binding && \
    go get github.com/martini-contrib/render && \
    go get github.com/martini-contrib/cors && \
    go get github.com/mikepb/go-serial

CMD ["/bin/bash"]

WORKDIR /work
