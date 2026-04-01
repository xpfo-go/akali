FROM golang:<xpfo{ .GoVersion }xpfo>-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/app .

FROM gcr.io/distroless/static-debian12

WORKDIR /app
COPY --from=builder /out/app /app/app
COPY config.yaml /app/config.yaml

EXPOSE 17878
ENTRYPOINT ["/app/app"]
