FROM golang:latest

WORKDIR /app

COPY . /app

CMD ["go", "run", "cmd/worker/opp_edge_worker.go"]
