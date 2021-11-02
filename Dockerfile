FROM golang:1.17-alpine AS build-env
WORKDIR /go/src/video
COPY go.mod go.sum ./
RUN go mod download
COPY . .

ARG TARGETARCH

RUN GOARCH=$TARGETARCH go build -o bin/client ./client/
RUN GOARCH=$TARGETARCH go build -o bin/server ./server/

FROM alpine:3.14
RUN apk add --no-cache ca-certificates iptables iproute2 ip6tables
COPY  --from=build-env /go/src/video/public /public
COPY --from=build-env /go/src/video/bin/* /usr/local/bin/