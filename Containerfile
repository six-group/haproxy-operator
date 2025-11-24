FROM golang:1.24.10-bookworm as builder

WORKDIR /workspace
COPY . /workspace

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags="-s -w" -o haproxy-operator main.go

RUN strip haproxy-operator

FROM gcr.io/distroless/static:nonroot
ENTRYPOINT ["/opt/go/haproxy-operator"]
WORKDIR /opt/go/
COPY --from=builder /workspace/haproxy-operator /opt/go/haproxy-operator
USER 1001:1001