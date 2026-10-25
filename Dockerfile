ARG GO_VERSION=1.27.1

FROM golang:${GO_VERSION} AS build
WORKDIR /src

# Copy the module files first so the dependency download is cached until
# go.mod or go.sum change.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# CGO off gives a fully static binary (lib/pq and godotenv are pure Go).
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/movies-api .

# Distroless static: no shell, no package manager, runs as a non-root user.
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/movies-api /movies-api
USER nonroot:nonroot
EXPOSE 8081
ENTRYPOINT ["/movies-api"]