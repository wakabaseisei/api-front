ARG GO_VERSION=1.24.0

FROM golang:${GO_VERSION}-bullseye AS builder

WORKDIR /app

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=bind,target=. \
    go build -o /bin/app

FROM gcr.io/distroless/base-debian12:nonroot

COPY --from=builder /bin/app /bin/app

ENTRYPOINT ["/bin/app"]
