FROM golang:1.21-bullseye as builder

WORKDIR /workspace
COPY . /workspace

RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -a -o haproxy-operator main.go

FROM scratch
WORKDIR /opt/go/
COPY --from=builder /workspace/haproxy-operator /opt/go/haproxy-operator
USER 1001:1001
ENTRYPOINT ["/opt/go/haproxy-operator"]
