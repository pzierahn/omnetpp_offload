# syntax=docker/dockerfile:1
FROM golang:latest AS builder

WORKDIR /install

COPY . /install
RUN rm -rf go.sum; \
    go get all
RUN go install cmd/worker/opp_offload_worker.go; \
    go install cmd/config/opp_offload_config.go; \
    go install cmd/broker/opp_offload_broker.go; \
    go install cmd/run/opp_offload_run.go; \
    go install cmd/stargate_client/stargate_client.go; \
    go install cmd/stargate_server/stargate_server.go

FROM pzierahn/omnetpp:6.0.1
WORKDIR /root

COPY --from=builder /go/bin/ /bin/
