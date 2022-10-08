FROM golang:alpine

WORKDIR /workspace

# this needs to be run from root of googlecloudplatform/guest-test-infra
COPY . .
RUN ls
RUN CGO_ENABLED=0 go build -o /agent ./agent/main.go
RUN CGO_ENABLED=0 go build -o /client ./client/main.go

RUN chmod +x /agent /client

FROM alpine

COPY --from=0 /agent /agent
COPY --from=0 /client /client
