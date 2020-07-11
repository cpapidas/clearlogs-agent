FROM golang:1.14 as builder
WORKDIR /src
COPY . .
# ldflags -w and -s result in smaller binary.
# -w Omits the DWARF symbol table.
# -s Omits the symbol table and debug information.
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -ldflags "-w -s -X main.theVersion=$(git describe --tags --always)" -o /bin/clagent ./cmd/clagent/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -ldflags "-w -s -X main.theVersion=$(git describe --tags --always)" -o /bin/demolog ./cmd/demolog/main.go

FROM alpine:latest
COPY --from=builder /bin/clagent .
COPY --from=builder /bin/demolog .
CMD ["/bin/sh", "-c", "./clagent -token=token -port=8090"]
CMD ["/bin/sh", "-c", "./demolog"]
