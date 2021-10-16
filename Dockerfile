FROM golang:1.16


WORKDIR /

ADD go.mod .
ADD go.sum .

RUN go mod download
ADD . .


EXPOSE 8000

RUN go build -o xapiens

CMD ["./xapiens"]