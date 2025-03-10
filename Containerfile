FROM golang:1.23.7-bullseye as builder
RUN apt-get update && apt-get install -y upx

WORKDIR /workspace
COPY . /workspace

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o haproxy-operator main.go && \
    upx -q haproxy-operator

FROM gcr.io/distroless/static:nonroot
ENTRYPOINT ["/opt/go/haproxy-operator"]
WORKDIR /opt/go/
COPY --from=builder /workspace/haproxy-operator /opt/go/haproxy-operator
USER 1001:1001