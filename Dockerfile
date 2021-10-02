FROM golang:1.17 AS builder

WORKDIR /go/src/
COPY . .

RUN go build main.go &&\
    ls

FROM golang:1.17
WORKDIR /go/src/
COPY --from=builder /go/src/main .

CMD ["./main"]
