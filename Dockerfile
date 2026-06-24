FROM golang:1.25.11-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /src

RUN addgroup -S nonroot && adduser -S nonroot -u 65532 -G nonroot

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -o /app/hookdrop .

FROM gcr.io/distroless/static-debian12:latest

COPY --from=builder /app/hookdrop /hookdrop
COPY --from=builder /etc/passwd /etc/passwd

USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/hookdrop"]
