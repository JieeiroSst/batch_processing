FROM golang:alpine

RUN mkdir /app

WORKDIR /app

ADD go.mod .
ADD go.sum .

RUN go mod download
ADD . .

EXPOSE 8000

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o /out/main ./main.go

ENTRYPOINT ["/out/main"]