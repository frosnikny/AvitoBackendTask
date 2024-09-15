FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd/app

RUN go build -o go-server

WORKDIR /app/cmd/migrate

RUN go build -o migrate

EXPOSE 8080

CMD ["sh", "-c", "./migrate && ./go-server"]
