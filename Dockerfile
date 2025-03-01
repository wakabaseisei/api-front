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
    go build -o /bin/app ./internal/cmd/app

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=bind,target=. \
    go build -o /bin/migrate-cli ./internal/cmd/cli

RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.2

FROM gcr.io/distroless/base-debian12:nonroot

COPY --from=builder /bin/app /bin/app
COPY --from=builder /bin/migrate-cli /bin/migrate-cli
COPY --from=builder /go/bin/migrate /bin/migrate
COPY db/migrations /db/migrations

ENTRYPOINT ["/bin/app"]
