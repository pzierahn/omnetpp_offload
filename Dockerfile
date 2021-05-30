# syntax=docker/dockerfile:1
FROM golang:latest AS builder

WORKDIR /install

COPY . /install
RUN go build cmd/worker/opp_edge_worker.go; \
    go build cmd/consumer/opp_edge_run.go

FROM ubuntu:latest
WORKDIR /root

COPY --from=builder /install/opp_edge* /bin
CMD ["opp_edge_worker"]
