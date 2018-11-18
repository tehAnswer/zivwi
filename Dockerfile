FROM golang:onbuild
ENV HOME /go/src/github.com/tehAnswer/zivwi
WORKDIR $HOME
COPY . ./
