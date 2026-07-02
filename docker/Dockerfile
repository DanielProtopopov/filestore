FROM golang:1.25-alpine as build

RUN apk add git --no-cache
WORKDIR /usr/local/app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN go build -ldflags "-s -w" -o filestore_bin ./cmd/api/

FROM alpine:3.18 as app
RUN apk add ca-certificates tzdata mysql-client --update && rm -rf /var/cache/apk/*

WORKDIR /usr/local/app
COPY --from=build /usr/local/app/filestore_bin bin/filestore
COPY ./data ./data

CMD ["/usr/local/app/bin/filestore"]
