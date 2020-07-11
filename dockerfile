FROM golang:1.13 as builder
WORKDIR /home/agent
COPY . ./
ARG version=dev
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -a -installsuffix cgo -ldflags "-X main.version=$version" -o agent ./cmd/agent/main.go

FROM alpine
COPY --from=builder /home/agent/agent .

ENTRYPOINT ["./agent"]