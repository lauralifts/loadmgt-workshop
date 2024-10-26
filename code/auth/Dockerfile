FROM golang:1.23
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

# config HTTP endpoint
EXPOSE 9010

RUN CGO_ENABLED=0 GOOS=linux go build -o /auth

CMD ["/auth"]