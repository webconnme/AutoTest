#!/bin/sh

rm -Rf pkg/*
rm -Rf src/github.com/

go get -tags zmq_4_x github.com/alecthomas/gozmq && \
    go get github.com/codeskyblue/go-sh && \
    go get github.com/satori/go.uuid && \
    go get github.com/go-martini/martini && \
    go get github.com/martini-contrib/binding && \
    go get github.com/martini-contrib/render && \
    go get github.com/martini-contrib/cors && \
    go get github.com/mikepb/go-serial


