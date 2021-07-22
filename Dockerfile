FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go log/*.go metrics/*.go ./

RUN go build -o /heroku-logs-exporter

ENTRYPOINT ["/heroku-logs-exporter"]
