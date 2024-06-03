FROM golang:1.22.3 as builder

RUN addgroup gouser && \
    adduser --ingroup gouser --shell /bin/false gouser && \
    cat /etc/passwd | grep gouser > /etc/passwd_gouser

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o server examples/middleware/main.go

FROM scratch
WORKDIR /app
COPY --from=builder /etc/passwd_gouser /etc/passwd
COPY --from=builder /app/server /app/.env ./

USER gouser

EXPOSE 8080

ENTRYPOINT ["./server"]
