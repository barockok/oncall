# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
	go mod download
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X github.com/barockok/oncall/engine-go/internal/version.version=${GIT_TAG}" -o /out/engine-go ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=builder /out/engine-go /app/engine-go
USER nonroot
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/app/engine-go"]