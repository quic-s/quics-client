FROM golang:alpine

WORKDIR /quics-client

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . .
RUN go build -o /qic ./cmd

WORKDIR /
COPY .env .
VOLUME [ "/dirs" ]

CMD [ "/qic", "start" ]