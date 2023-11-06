FROM golang:alpine

WORKDIR /quics-client
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o qic ./cmd

EXPOSE 6121/udp

ENV PATH="/quics-client:${PATH}"

RUN chmod 777 /root
RUN mkdir /dirs
RUN chmod 777 /dirs

VOLUME [ "/dirs" ]

CMD [ "qic", "start" ]