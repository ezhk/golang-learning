FROM golang:latest as builder

WORKDIR /app
COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go get -d -v
RUN go get -u github.com/pressly/goose/cmd/goose

RUN go build -tags=server,sender,scheduler -ldflags="-w -s" -o calendar .


# New stage
FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/calendar .
COPY --from=builder /go/bin/goose .

COPY migrations /app/migrations/
COPY deployments/configs /app/configs/
