FROM golang:latest
RUN mkdir /app
WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download

COPY . /app/
RUN go build -o main ./cmd/goupload
CMD ["/app/main"]
