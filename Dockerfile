FROM golang:alpine

ENV CGO_ENABLED=0
WORKDIR /workspace

# this needs to be run from root of googlecloudplatform/guest-test-infra
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /agent ./agent/main.go
RUN go build -o /client ./client/main.go

RUN chmod +x /agent /client

FROM alpine

COPY --from=0 /agent /agent
COPY --from=0 /client /client
