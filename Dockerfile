FROM golang

WORKDIR /usr/src/app

COPY go.mod ./
RUN go mod download && go mod verify

COPY . .

RUN mkdir -p /usr/local/bin/go-imdg
RUN go build -v -o /usr/local/bin/go-imdg ./...

CMD ["/usr/local/bin/go-imdg/node"]