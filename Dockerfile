FROM golang:latest as builder

WORKDIR /app
COPY . .

RUN go build -o cold main.go

FROM ubuntu:latest

COPY --from=builder /app/cold /usr/local/bin/
COPY cold.db /root/cold.db

EXPOSE 8080
ENTRYPOINT ["cold"]
