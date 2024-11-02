FROM golang:1.23.1 AS build-stage
WORKDIR /app/src
COPY . ./

RUN go mod download

WORKDIR /app/src/server
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/fci

FROM debian:bookworm-slim AS builder

FROM build-stage AS run-test-stage
RUN go test -v ./...

FROM gcr.io/distroless/base-debian12 AS build-release-stage
WORKDIR /app/bin
COPY --from=build-stage /app/bin/fci ./fci

WORKDIR /opt/fci
COPY --from=build-stage /app/src/config.yml ./config.yml

COPY --from=builder /bin/sh /bin/sh
COPY --from=builder /bin/mkdir /bin/mkdir
COPY --from=builder /bin/cat /bin/cat

EXPOSE 8888

ENTRYPOINT ["/app/bin/fci"]