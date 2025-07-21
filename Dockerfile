FROM golang:latest

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./

COPY . .

RUN go mod download

EXPOSE 3000

CMD ["air", "-c", ".air.toml"]
