# Builder

# ARG GITHUB_PATH=github.com/aperg/my-api

FROM golang:1.20-alpine AS builder
RUN apk add --update make git protoc protobuf protobuf-dev curl
COPY . /home/${GITHUB_PATH}
WORKDIR /home/${GITHUB_PATH}
# RUN make deps-go
RUN make build-go

# gRPC Server

FROM alpine:latest as server
LABEL org.opencontainers.image.source https://${GITHUB_PATH}
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /home/${GITHUB_PATH}/bin/grpc-server .
COPY --from=builder /home/${GITHUB_PATH}/config.yml .

RUN chown root:root grpc-server

EXPOSE 50051
EXPOSE 8080
EXPOSE 9100

CMD ["./grpc-server"]
