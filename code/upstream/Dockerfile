FROM golang:1.23
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY .env ./

EXPOSE 9092
EXPOSE 9093

RUN CGO_ENABLED=0 GOOS=linux go build -o /upstream

CMD ["/upstream"]