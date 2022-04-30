FROM golang:1.16

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./cmd/honeypot-ingestion
RUN go build ./cmd/honeypot-ingestion

CMD ./honeypot-ingestion