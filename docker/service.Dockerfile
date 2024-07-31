FROM golang:1.22-alpine as builder

WORKDIR /app

COPY go.mod go.sum* ./

RUN go mod download

COPY . .

RUN go build -o main cmd/main.go

EXPOSE 5001

CMD sh -c "sleep 25 && ./main"