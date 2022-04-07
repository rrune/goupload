FROM golang:latest
RUN mkdir /app
WORKDIR /app
RUN mkdir ./src
WORKDIR /app/src
ADD ./src /app/src
RUN go mod download
RUN go build -o main ./cmd
CMD ["/app/src/main"]