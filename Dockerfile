# syntax=docker/dockerfile:1
FROM golang:latest AS builder

WORKDIR /install

COPY . /install
RUN rm -rf go.sum; \
    go get all
RUN go build cmd/worker/opp_edge_worker.go; \
    go build cmd/config/opp_edge_config.go; \
    go build cmd/broker/opp_edge_broker.go

FROM pzierahn/omnetpp
WORKDIR /root

COPY --from=builder /install/opp_edge* /bin/
