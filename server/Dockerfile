FROM golang:1.25 AS builder

WORKDIR /app

RUN touch db.sqlite

RUN go install github.com/go-task/task/v3/cmd/task@latest

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG LDFLAGS

RUN go build -ldflags="$LDFLAGS" .

FROM cloudflare/cloudflared:latest

WORKDIR /app

COPY --from=builder --chown=1000:1000 /app /app

EXPOSE 8090

ENTRYPOINT ["/app/cloudflared-web-gui"]
