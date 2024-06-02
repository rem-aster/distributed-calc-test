FROM golang:1.19-alpine
WORKDIR /app
COPY . .
RUN go build -o orch .
EXPOSE 8080
ENTRYPOINT ["./orch"]
