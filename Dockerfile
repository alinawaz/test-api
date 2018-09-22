FROM golang:1.8

ENV GOPATH $HOME/go
ENV PATH $HOME/go/bin:$PATH

ADD . /go/src/test-api
WORKDIR /go/src/test-api

RUN go get test-api
RUN go install

ENTRYPOINT ["/go/bin/test-api"]

EXPOSE 8080