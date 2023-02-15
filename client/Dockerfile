FROM mcr.microsoft.com/oss/go/microsoft/golang:1.19 AS builder

WORKDIR /go-grpc/client

COPY . .

RUN CGO_ENABLED=1 go build -v -o /usr/local/bin/grpcClient -gcflags="-dwarflocationlists=true"  main.go

FROM mcr.microsoft.com/oss/mirror/docker.io/library/ubuntu:20.04
COPY --from=builder /usr/local/bin/grpcClient \
  /usr/bin/grpcClient

COPY --from=builder /go-grpc/client \
  /usr/local/script

RUN chmod +x /usr/bin/grpcClient
RUN chmod +rx /usr/local/script/generate_certs.sh

WORKDIR /usr/local/script
RUN apt-get update
RUN apt-get install -y openssl
RUN apt-get install -y ca-certificates
RUN apt-get upgrade -y
RUN ./generate_certs.sh

ENTRYPOINT ["/usr/bin/grpcClient", "start"]
