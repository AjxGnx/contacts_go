FROM golang:1.20-alpine AS build

WORKDIR /go/src/github.com/AjxGnx/contacts_go

COPY . .

RUN go mod download

RUN go build -o app ./cmd/main.go

FROM alpine:3.12

WORKDIR /usr/

COPY --from=build /go/src/github.com/AjxGnx/contacts_go .

CMD /usr/app