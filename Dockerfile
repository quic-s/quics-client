FROM golang:alpine

WORKDIR /quics-client
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /bin/qic ./cmd

EXPOSE 6121/udp

WORKDIR /
COPY .env ./bin
VOLUME [ "/dirs" ]

CMD [ "qic", "start" ]