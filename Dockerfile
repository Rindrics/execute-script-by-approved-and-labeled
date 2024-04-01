FROM golang:1.22 AS builder

WORKDIR /dist

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .


FROM alpine:3.19

LABEL org.opencontainers.image.source="https://github.com/Rindrics/execute-scripts-github-flow"

RUN apk update && apk --no-cache add ca-certificates

COPY --from=builder /dist/app /app

CMD ["/app"]