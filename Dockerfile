FROM golang:1.27.1

RUN go install github.com/air-verse/air@v1.66.0

WORKDIR /app
CMD ["air", "-c", ".air.toml"]
